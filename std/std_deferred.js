import { __std as def } from '__std_Deferred_generated';
import { __std as m } from '__std_generated';
import flatbuffers from 'flatbuffers';

const deferreds = {};

function recv(buf) {
  const data = new flatbuffers.ByteBuffer(new Uint8Array(buf));
  const reso = def.Fulfilment.getRootAsFulfilment(data);
  const ser = reso.serial().toFloat64();
  let callback;
  let value;
  switch (reso.valueType()) {
  case def.FulfilmentValue.Data: {
    ({ data: callback } = deferreds[ser]);
    const val = new def.Data();
    reso.value(val);
    value = val.bytes();
    break;
  }
  case def.FulfuiilmentValue.Error: {
    ({ error: callback } = deferreds[ser]);
    const err = new def.Error();
    reso.value(err);
    value = new Error(err.error());
    break;
  }
  case def.FulfilmentValue.EndOfStream:
    ({ end: callback } = deferreds[ser]);
    break;
  default:
    throw new Error('Unknown message recieved from runtime');
  }

  if (callback === undefined) {
    // for now, drop it, on the basis that cancelling may be racey.
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
function requestAsPromise(fn) {
  const buf = fn();
  const data = new flatbuffers.ByteBuffer(new Uint8Array(buf));
  const resp = def.DeferredResponse.getRootAsDeferredResponse(data);
  switch (resp.retvalType()) {
  case def.DeferredRetval.Error: {
    const err = new def.Error();
    resp.retval(err);
    return Promise.reject(new Error(err.error()));
  }
  case def.DeferredRetval.Deferred: {
    const defer = new def.Deferred();
    resp.retval(defer);
    const ser = defer.serial().toFloat64();
    return new Promise((resolve, reject) => {
      function removeThenCall(v) {
        delete deferreds[ser];
        return this(v);
      }
      registerCallbacks(ser,
        removeThenCall.bind(resolve),
        removeThenCall.bind(reject),
        removeThenCall.bind(panic('Unexpected EndOfStream for promisified deferred')));
    });
  }
  default:
    return Promise.reject(new Error('Failed to decode response from request'));
  }
}

// TODO
function cancel(serial) {
  const builder = new flatbuffers.Builder(512);
  def.ReadArgs.startCancelArgs(builder);
  def.ReadArgs.addSerial(builder, serial);
  const argsOffset = def.CancelArgs.endCancelArgs(builder);

  m.Message.startMessage(builder);
  m.Message.addArgsType(builder, m.Args.CancelArgs);
  m.Message.addArgs(builder, argsOffset);
  const messageOffset = m.Message.endMessage(builder);
  builder.finish(messageOffset);
  return V8Worker2.send(builder.asArrayBuffer());
}

V8Worker2.recv(recv);

export {
  requestAsPromise,
  cancel,
};
