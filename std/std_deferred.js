import { __std as rpc } from '__std_RPC_generated';
import flatbuffers from 'flatbuffers';

var deferreds = {};

function recv(buf) {
    const data = new flatbuffers.ByteBuffer(new Uint8Array(buf));
    const reso = rpc.Resolution.getRootAsResolution(data);
    const ser = reso.serial().toFloat64();
    var callback, value;
    switch (reso.valueType()) {
    case rpc.ResolutionValue.Data:
        ({data: callback} = deferreds[ser]);
        const val = new rpc.Data();
        reso.value(val);
        value = val.bytes();
        break
    case rpc.ResolutionValue.Error:
        ({error: callback} = deferreds[ser]);
        const err = new rpc.Error()
        reso.value(err);
        value = new Error(err.error());
        break
    case rpc.ResolutionValue.EndOfStream:
        ({end: callback} = deferreds[ser]);
        break
    default:
        throw new Error('Unknown message recieved from runtime');
    }

    if (callback === undefined) {
        // for now, drop it, on the basis that cancelling may be racey.
        return
    }
    callback(value);
}

// registerCallbacks records callbacks for the three possible outcomes
// of a deferred, which is identified by a serial number returned from
// Go.
function registerCallbacks(serial, onData, onError, onEnd) {
    deferreds[serial] = {data: onData, error: onError, end: onEnd}
}

function panic(msg) {
    return function() {
        throw new Error(msg);
    }
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
    const resp = rpc.Response.getRootAsResponse(data);
    switch (resp.retvalType()) {
    case rpc.Retval.Error:
        const err = new rpc.Error();
        resp.retval(err);
        return Promise.reject(new Error(err.error()))
    case rpc.Retval.Deferred:
        const defer = new rpc.Deferred();
        resp.retval(defer);
        const ser = defer.serial().toFloat64(); 
        return new Promise(function(resolve, reject) {
            function removeThenCall(v) {
                delete deferreds[ser];
                return this(v);
            }
            registerCallbacks(ser,
                              removeThenCall.bind(resolve),
                              removeThenCall.bind(reject),
                              removeThenCall.bind(panic('Unexpected EndOfStream for promisified deferred')));
        });
    default:
        return Promise.reject(new Error('Failed to decode response from request'));
    }
}

// TODO
function cancel(serial) {
}

V8Worker2.recv(recv);

export {
    requestAsPromise,
};
