export type Data = Uint8Array | string;
export type Transform = (x: Data) => Data;

export const ident : Transform = (x: Data): Data => x;

function uint8ToUint16Array(bytes: Uint8Array): Uint16Array {
  return new Uint16Array(bytes.buffer, bytes.byteOffset, bytes.byteLength / 2);
}

// From the runtime, we get various encodings back depending on the
// type in the schema, some of which will be post-processed by
// flatbuffers.
export const stringFromUTF16Bytes = (bytes: Uint8Array): string => String.fromCodePoint(...uint8ToUint16Array(bytes));
export const stringFromUTF8Bytes = (bytes: Uint8Array): string => String.fromCodePoint(...bytes);
export const valueFromUTF8Bytes = (bytes: Uint8Array): any => JSON.parse(stringFromUTF8Bytes(bytes));
export const valueFromUTF16Bytes = (bytes: Uint8Array): any => JSON.parse(stringFromUTF16Bytes(bytes));
