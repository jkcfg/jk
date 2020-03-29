import * as std from '../index';
import { WriteOptions } from '../write';
import { ValidateFn } from './validate';
import { normaliseResult, formatError } from '../validation';

/* eslint @typescript-eslint/explicit-function-return-type: "off" */

/**
 * File is the basic unit of input for generate; it represents a
 * configuration item to output to a file.
 */
export interface File {
  path: string;
  value: any | Promise<any>;
  format?: std.Format;
  validate?: ValidateFn;
}

/*
 * GenerateParams types the optional arguments to generate.
 */
export interface GenerateParams {
  stdout?: boolean;
  overwrite?: std.Overwrite;
  writeFile?: (v: any, p: string, o?: WriteOptions) => void;
}

const helpMsg = `
To use generate, export a default value with the list of files to generate:

  export default [
    { path: 'file1.yaml', value: value1 },
    { path: 'file2.yaml', value: [v0, v1, v2], format: std.Format.YAMLStream },
    ...
  ];

Notes:

- The default export can also be a promise to such a array.
- Optional parameters are the same as std.write().`;

function error(msg: string): void {
  std.log(msg);
}

function help(): void {
  std.log(helpMsg);
}

/**
 * Calculate the modulus of two numbers
 * @param {number} x
 * @param {number} y
 * @returns {number} res
 * @private
 */
function mod(x: number, y: number): number {
  if (y > 0) {
    // We don't use JavaScript's modulo operator here as this doesn't work
    // correctly for x < 0 and x === 0
    // see https://en.wikipedia.org/wiki/Modulo_operation
    return x - y * Math.floor(x / y);
  }
  if (y === 0) {
    return x;
  }
  // TODO: implement mod for a negative divisor
  throw new Error('cannot calculate mod for a negative divisor');
}

// https://stackoverflow.com/questions/13627308/add-st-nd-rd-and-th-ordinal-suffix-to-a-number/13627586
// We don't use the modulo operator here as it would be interpreted by Sprintf.
const nth = (n: number): string => {
  const s = ['th', 'st', 'nd', 'rd'];

  const v = mod(n, 100);
  return n + (s[mod(v - 20, 10)] || s[v] || s[0]);
};

function extension(path: string): string {
  return path.split('.').pop();
}

function formatFromPath(path: string): std.Format {
  switch (extension(path)) {
  case 'yaml':
  case 'yml':
    return std.Format.YAML;
  case 'json':
    return std.Format.JSON;
  case 'hcl':
  case 'tf':
    return std.Format.HCL;
  default:
    return std.Format.JSON;
  }
}

const isString = (s: any): boolean => typeof s === 'string' || s instanceof String;

// represents a file spec that has its promise resolved, if necessary
interface RealisedFile {
  value: any;
  format?: std.Format;
  file?: string;
  path?: string;
}

// Compute the output format of a file spec.
function fileFormat(o: RealisedFile): std.Format {
  let { path, format, value } = o;

  if (format === undefined || format === std.Format.FromExtension) {
    if (isString(value)) {
      format = std.Format.Raw;
    } else {
      format = formatFromPath(path);
    }
  }

  return format;
}

const formatName = (f: std.Format): string => std.Format[f];

function usedFormats(summary: number[]): string[] {
  const augmented = summary.map((n, i) => ({ format: formatName(i), n }));
  return augmented.reduce((formats, desc) => {
    if (desc.n > 0) {
      formats.push(desc.format);
    }
    return formats;
  }, []);
}

function validateFormat(files: RealisedFile[], params: GenerateParams) {
  /* we have an array */
  if (!Array.isArray(files)) {
    error('default value is not an array');
    return { valid: false, showHelp: true };
  }

  /* an array with each element a { path, value } object */
  let valid = true;
  files.forEach((e, i) => {
    /* 'file' is the old 'path' property name. Fixup things */
    if (e.file !== undefined) {
      e.path = e.file;
    }

    ['path', 'value'].forEach((prop) => {
      if (!Object.prototype.hasOwnProperty.call(e, prop)) {
        error(`${nth(i + 1)} element does not have a '${prop}' property`);
        valid = false;
      }
    });
  });

  return { valid, showHelp: !valid };
}

function assembleForStdout(values: RealisedFile[]) {
  // When writing to stdout, we need to
  //  1. make sure everything is a mutually compatible format (e.g.,
  //     YAMLs and YAMLStreams, but not JSON with YAMLStreams);
  //  2. collate the values to be printed as a streamable array, which
  //     means anything that's _already_ a stream value is inlined.
  let stdoutFormat;
  let stream: any[] = [];
  const formatsSeen = Array(Object.keys(std.Format).length).fill(0);

  for (const v of values) {
    const format = fileFormat(v);
    formatsSeen[format] += 1;

    switch (format) {
    case std.Format.YAML:
    case std.Format.YAMLStream:
      if (stdoutFormat !== undefined && stdoutFormat !== std.Format.YAMLStream) {
        error(`stdout requires compatible formats, but have seen ${usedFormats(formatsSeen).join(',')}`);
        return { valid: false }
      }
      stdoutFormat = std.Format.YAMLStream;
      stream = (format === std.Format.YAML) ?
        stream.push(v.value) && stream : stream.concat(v.value);
      break;
    case std.Format.JSON:
    case std.Format.JSONStream:
      if (stdoutFormat !== undefined && stdoutFormat !== std.Format.JSONStream) {
        error(`stdout requires compatible formats, but have seen ${usedFormats(formatsSeen).join(',')}`);
        return { valid: false }
      }
      stdoutFormat = std.Format.JSONStream;
      stream = (format === std.Format.JSON) ?
        stream.push(v.value) && stream : stream.concat(v.value);
      break;
    default:
      // for anything else, only one value is allowed; therefore keep
      // the value as it is, but check that this is the only value.
      if (stdoutFormat !== undefined) {
        error(`stdout requires compatible formats, but have seen ${usedFormats(formatsSeen).join(',')}`);
        return { valid: false }
      }
      stdoutFormat = format;
      stream = v.value;
      break;
    }
  }

  return { valid: true, stdoutFormat, stream }
}

type GenerateArg = File[] | Promise<File[]> | (() => File[]);

/**
 * generate is the entry point for the module; it accepts the input
 * configuration, and outputs validated values to the files as
 * specified.
 */
export function generate(definition: GenerateArg, params: GenerateParams) {
  /*
   * The default export can be:
   *  1. an array of { path, value } objects,
   *  2. a promise to such an array,
   *  3. a function evaluating to either 1. or 2.
   */
  let inputs: Promise<File[]>;
  if (typeof definition === 'function') {
    inputs = Promise.resolve(definition());
  } else {
    inputs = Promise.resolve(definition);
  }

  const { stdout = false, overwrite = false, writeFile = std.write } = params;

  return async function() {
    const files = await inputs;
    const resolved = await Promise.all(files.map(f => f.value));
    resolved.forEach((v, i) => {
      // eslint-disable-next-line no-param-reassign
      files[i].value = v;
    });

    // check the format of the generated value
    const { valid: formatValid, showHelp } = validateFormat(files, params);
    if (showHelp) {
      help();
    }
    if (!formatValid) {
      throw new Error('jk-internal-skip: format invalid');
    }

    let results = await Promise.all(files.map(({ path, value, validate = (() => 'ok') }) => {
      return Promise.resolve(validate(value))
        .then(r => ({ path, result: normaliseResult(r) }));
    }));

    let valuesValid = true;
    results.forEach(({ path, result }) => {
      if (result !== 'ok') {
        result.forEach(err => error(formatError(path, err)));
        valuesValid = false;
      }
    });

    if (!valuesValid) {
      throw new Error('jk-internal-skip: values failed validation');
    }

    if (stdout) {
      const { valid, stdoutFormat, stream } = assembleForStdout(files);
      if (!valid) {
        throw new Error('jk-internal-skip: validation failed');
      }
      std.write(stream, std.stdout, { format: stdoutFormat });
    } else {
      for (const o of files) {
        const { path, value, ...args } = o;
        writeFile(value, path, { overwrite, ...args });
      }
    }
  }();
}
