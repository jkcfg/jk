import { requestAsPromise } from 'std_deferred';
import flatbuffers from 'flatbuffers';
import { __std as r } from '__std_Read_generated';
import { __std as m } from '__std_generated';

// read requests a URL and returns a promise that will be resolved
// with the contents at the URL, or rejected.
function read(url) {
    const builder = new flatbuffers.Builder(512);
    const urlOffset = builder.createString(url);
    r.ReadArgs.startReadArgs(builder);
    r.ReadArgs.addUrl(builder, urlOffset);
    const argsOffset = r.ReadArgs.endReadArgs(builder);

    m.Message.startMessage(builder);
    m.Message.addArgsType(builder, m.Args.ReadArgs);
    m.Message.addArgs(builder, argsOffset);
    const messageOffset = m.Message.endMessage(builder);
    builder.finish(messageOffset);
    return requestAsPromise(() => {
        return V8Worker2.send(builder.asArrayBuffer());
    });
}

// TODO(michael): read with format?

export {
    read
}
