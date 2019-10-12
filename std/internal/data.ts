export type Transform = (x: Uint8Array) => any;

export const ident : Transform = (x: Uint8Array): any => x;

// There's a practical limit to the number of arguments that can be
// given to a function, either via apply or spread syntax. This
// function chunks the conversion using a conservative number of
// arguments to String.fromCharCode; the exact limit for v8 appears to
// be 65535, but evidence elsewhere suggests a smaller value is
// better:
// https://github.com/google/closure-library/commit/da353e0265ea32583ea1db9e7520dce5cceb6f6a
function stringFromUTF16CharCodes(bytes: ArrayBuffer, byteOffset: number, byteLength: number): string {
  const chunk = 8192 * 2;
  let result = '';
  for (let i = 0; i < byteLength; i += chunk) {
    const len = ((i + chunk) > byteLength) ? (byteLength - i)/2 : chunk / 2;
    const codes = new Uint16Array(bytes, byteOffset + i, len);
    result += String.fromCharCode.apply(null, codes);
  };
  return result;
}

// From the runtime, we get various encodings back depending on the
// type in the schema, some of which will be post-processed by
// flatbuffers.
export const stringFromUTF16Bytes = (bytes: Uint8Array): string => stringFromUTF16CharCodes(bytes.buffer, bytes.byteOffset, bytes.byteLength);
export const stringFromUTF8Bytes = (bytes: Uint8Array): string => String.fromCharCode.apply(null, bytes);
export const valueFromUTF8Bytes = (bytes: Uint8Array): any => JSON.parse(stringFromUTF8Bytes(bytes));
export const valueFromUTF16Bytes = (bytes: Uint8Array): any => JSON.parse(stringFromUTF16Bytes(bytes));
