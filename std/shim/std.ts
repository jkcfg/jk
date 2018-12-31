function print(...args: any[]): void { };
function log(value: any): void { };

enum Format {
    JSON = 1,
    YAML,
    Raw,
}

interface WriteOptions {
  format?: Format,
  indent?: number,
  override?: boolean,
}

function write(value: any, path: string, options?: WriteOptions): void { };

enum Encoding {
  Bytes,
  String,
  JSON,
}

interface ReadOptions {
  format: Format,
  encoding: Encoding,
}

function read(path: string, options?: ReadOptions): Promise<any> { return Promise.resolve({}) };

export default {
  print, log, Format, write, read,
};
