function stringToArrayBuffer(s) {
  const buf = new ArrayBuffer(s.length * 2);
  const view = new Uint16Array(buf);
  for (let i = 0, l = s.length; i < l; i ++) {
    view[i] = s.charCodeAt(i);
  }
  return buf;
}

function write(value) {
    const json = JSON.stringify(value);
    const buf = stringToArrayBuffer(json);
    V8Worker2.send(buf);
}

export default write;
