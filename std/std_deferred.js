import { __std } from '__std_generated';
import flatbuffers from 'flatbuffers';

const deferreds = {};

function recv(buf) {
  const data = new flatbuffers.ByteBuffer(new Uint8Array(buf));
  const reso = __std.Fulfilment.getRootAsFulfilment(data);
  const ser = reso.serial().toFloat64();
  let callback;
  let value;
  switch (reso.valueType()) {
  case __std.FulfilmentValue.Data: {
    ({ data: callback } = deferreds[ser]);
    const val = new __std.Data();
    reso.value(val);
    value = val.bytesArray();
    break;
  }
  case __std.FulfilmentValue.Error: {
    ({ error: callback } = deferreds[ser]);
    const err = new __std.Error();
    reso.value(err);
    value = new Error(err.message());
    break;
  }
  case __std.FulfilmentValue.EndOfStream:
    ({ end: callback } = deferreds[ser]);
    break;
  default:
    throw new Error('Unknown message recieved from runtime');
  }

  if (callback === undefined) {
    // for now, drop it, on the presumption that cancelation is
    // underway.
    return;
  }
  callback(value);
}

// registerCallbacks records callbacks for the three possible outcomes
// of a deferred, which is identified by a serial number returned from
// Go.
function registerCallbacks(serial, onData, onError, onEnd) {
  deferreds[serial] = { data: onData, error: onError, end: onEnd };
}

function panic(msg) {
  return () => { throw new Error(msg); };
}

// TODO(michael): consider factoring the V8Worker.send(...) calls into
// a function here as well.

// requestAsPromise performs the request given, and wraps the deferred
// result in a promise. If the request provokes an error, the Promise
// is rejected immediately; otherwise, the Promise will later be
// resolved or rejected depending on what is sent by the runtime.
function requestAsPromise(req, tx) {
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
      function removeThenCall(v) {
        delete deferreds[ser];
        return this(v);
      }
      function rejectWithStack(err) {
        /* eslint-disable no-param-reassign */
        err.stack += stackCapture.stack.substring(stackCapture.stack.indexOf('\n'));
        /* eslint-enable */
        reject(err);
      }
      const ondata = bytes => resolve(tx(bytes));
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
function cancel(serial) {
  const builder = new flatbuffers.Builder(512);
  __std.ReadArgs.startCancelArgs(builder);
  __std.ReadArgs.addSerial(builder, serial);
  const argsOffset = __std.CancelArgs.endCancelArgs(builder);

  __std.Message.startMessage(builder);
  __std.Message.addArgsType(builder, __std.Args.CancelArgs);
  __std.Message.addArgs(builder, argsOffset);
  const messageOffset = __std.Message.endMessage(builder);
  builder.finish(messageOffset);
  return V8Worker2.send(builder.asArrayBuffer());
}

V8Worker2.recv(recv);

export {
  requestAsPromise,
  cancel,
};
