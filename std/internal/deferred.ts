import { __std } from './__std_generated';
import { Transform } from './data';
import { flatbuffers } from './flatbuffers';

type Callback = (bytes: Uint8Array) => void;
type ErrorCallback = (err: Error) => void;

type Serial = number;

interface Deferred {
  data: Callback;
  end: Callback;
  error: ErrorCallback;
}

const deferreds: Map<Serial, Deferred> = new Map();

function recv(buf: ArrayBuffer): void {
  const data = new flatbuffers.ByteBuffer(new Uint8Array(buf));
  const reso = __std.Fulfilment.getRootAsFulfilment(data);
  const ser = reso.serial().toFloat64();
  let callback;
  let value;
  switch (reso.valueType()) {
  case __std.FulfilmentValue.Data: {
    ({ data: callback } = deferreds.get(ser));
    const val = new __std.Data();
    reso.value(val);
    value = val.bytesArray();
    break;
  }
  case __std.FulfilmentValue.Error: {
    const { error: errorCallback } = deferreds.get(ser);
    if (errorCallback === undefined) {
      return;
    }
    const err = new __std.Error();
    reso.value(err);
    const error = new Error(err.message());
    errorCallback(error);
    return;
  }
  case __std.FulfilmentValue.EndOfStream:
    ({ end: callback } = deferreds.get(ser));
    break;
  default:
    throw new Error('Unknown message received from runtime');
  }

  if (callback === undefined) {
    // for now, drop it, on the presumption that cancelation is
    // underway.
    return;
  }
  callback(value);
}

// `recv` is our handler for bytes sent ad-hoc from the runtime.
V8Worker2.recv(recv);

// registerCallbacks records callbacks for the three possible outcomes
// of a deferred, which is identified by a serial number returned from
// Go.
function registerCallbacks(serial: Serial, onData: Callback, onError: ErrorCallback, onEnd: Callback): void {
  deferreds.set(serial, { data: onData, error: onError, end: onEnd });
}

function panic(msg: string): ((_: ArrayBuffer) => never) {
  return () => { throw new Error(msg); };
}

// sendRequest sends an ArrayBuffer to the jk runtime and returns the
// ArrayBuffer response.
function sendRequest(buf: ArrayBuffer): ArrayBuffer {
  return V8Worker2.send(buf);
}

// requestAsPromise performs the request given, and wraps the deferred
// result in a promise. If the request provokes an error, the Promise
// is rejected immediately; otherwise, the Promise will later be
// resolved or rejected depending on what is sent by the runtime.
function requestAsPromise(req: () => ArrayBuffer, tx: Transform): Promise<any> {
  const buf = req();
  const data = new flatbuffers.ByteBuffer(new Uint8Array(buf));
  const resp = __std.DeferredResponse.getRootAsDeferredResponse(data);
  switch (resp.retvalType()) {
  case __std.DeferredRetval.Error: {
    const err = new __std.Error();
    resp.retval(err);
    return Promise.reject(new Error(err.message()));
  }
  case __std.DeferredRetval.Deferred: {
    const stackCapture = new Error();
    const defer = new __std.Deferred();
    resp.retval(defer);
    const ser = defer.serial().toFloat64();
    return new Promise((resolve, reject) => {
      function removeThenCall<V>(v: V) {
        deferreds.delete(ser);
        return this(v);
      }
      function rejectWithStack(err: Error) {
        /* eslint-disable no-param-reassign */
        err.stack += stackCapture.stack.substring(stackCapture.stack.indexOf('\n'));
        /* eslint-enable */
        reject(err);
      }
      const ondata = (bytes: Uint8Array) => resolve(tx(bytes));
      registerCallbacks(
        ser,
        removeThenCall.bind(ondata),
        removeThenCall.bind(rejectWithStack),
        removeThenCall.bind(panic('Unexpected EndOfStream for promisified deferred')),
      );
    });
  }
  default:
    return Promise.reject(new Error('Failed to decode response from request'));
  }
}

// TODO
function cancel(serial: Serial): ArrayBuffer {
  const builder = new flatbuffers.Builder(512);
  __std.CancelArgs.startCancelArgs(builder);
  __std.CancelArgs.addSerial(builder, builder.createLong(serial, 0));
  const argsOffset = __std.CancelArgs.endCancelArgs(builder);

  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.CancelArgs);
  __std.Message.addArgs(builder, argsOffset);
  const messageOffset = __std.Message.endMessage(builder);
  builder.finish(messageOffset);
  return sendRequest(builder.asArrayBuffer());
}

export {
  requestAsPromise,
  sendRequest,
  cancel,
};
