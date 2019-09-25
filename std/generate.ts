import * as std from './index';

/* eslint @typescript-eslint/explicit-function-return-type: "off" */

type ValidateFn = (value: any) => string;

export interface Value {
  path: string;
  value: any | Promise<any>;
  format?: std.Format;
  validate?: ValidateFn;
}

export interface GenerateParams {
  stdout?: boolean;
  overwrite?: std.Overwrite;
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
  std.log(`error: ${msg}`);
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

// represents a value that has has any promises resolved
interface RealisedValue {
  value: any;
  format?: std.Format;
  file?: string;
  path?: string;
}

// Compute the output format of a value.
function valueFormat(o: RealisedValue): std.Format {
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

function formatSummary(values: RealisedValue[]): number[] {
  const formats = Array(Object.keys(std.Format).length).fill(0);

  values.forEach((e): void => {
    formats[valueFormat(e)] += 1;
  });

  return formats;
}

const formatNames = [
  'FromExtension',
  'JSON',
  'YAML',
  'Raw',
  'YAMLStream',
  'JSONStream',
  'HCL',
];

const formatName = (f: std.Format): string => formatNames[f];

function usedFormats(summary: number[]): string[] {
  const augmented = summary.map((n, i) => ({ format: formatName(i), n }));
  return augmented.reduce((formats, desc) => {
    if (desc.n > 0) {
      formats.push(desc.format);
    }
    return formats;
  }, []);
}

function validate(values: RealisedValue[], params: GenerateParams) {
  /* we have an array */
  if (!Array.isArray(values)) {
    error('default value is not an array');
    return { valid: false, showHelp: true };
  }

  /* an array with each element a { path, value } object */
  let valid = true;
  values.forEach((e, i) => {
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

  if (!valid) {
    return { valid, showHelp: true };
  }

  /* when outputting to stdout, ensure that: */
  let stdoutFormat;
  if (params.stdout === true) {
    /* there's a single output format defined */
    const summary = formatSummary(values);
    const formats = usedFormats(summary);
    if (formats.length > 1) {
      error(`stdout output requires using a single format but got: ${formats.join(',')}`);
      return { valid: false, showHelp: false };
    }

    /*
     * If we have more than one file to generate, make sure it's either JSON or
     * YAML so we can output a stream of documents.
     */
    if (values.length > 1 && formats[0] !== 'JSON' && formats[0] !== 'YAML') {
      error(`stdout output for multiple files requires either JSON or YAML format but got: ${formats[0]}`);
      return { valid: false, showHelp: false };
    }

    if (values.length > 1) {
      if (formats[0] === 'JSON') {
        stdoutFormat = std.Format.JSONStream;
      } else if (formats[0] === 'YAML') {
        stdoutFormat = std.Format.YAMLStream;
      }
    } else {
      stdoutFormat = valueFormat(values[0]);
    }
  }

  return { valid: true, stdoutFormat, showHelp: false };
}

type ValueArg = Value[] | Promise<Value[]> | (() => Value[]);

function generate(definition: ValueArg, params: GenerateParams) {
  /*
   * The default export can be:
   *  1. an array of { path, value } objects,
   *  2. a promise to such an array,
   *  3. a function evaluating to either 1. or 2.
   */
  let inputs: Promise<Value[]>;
  if (typeof definition === 'function') {
    inputs = Promise.resolve(definition());
  } else {
    inputs = Promise.resolve(definition);
  }

  const { stdout = false, overwrite = false } = params;

  inputs.then((files) => {
    /* values can be promises as well */
    const justValues = files.map(f => f.value);
    Promise.all(justValues).then((resolved) => {
      resolved.forEach((v, i) => {
        // eslint-disable-next-line no-param-reassign
        files[i].value = v;
      });

      const { valid, stdoutFormat, showHelp } = validate(files, params);
      if (showHelp) {
        help();
      }
      if (!valid) {
        throw new Error('jk-internal-skip: validation failed');
      }

      if (stdout) {
        if (files.length > 1) {
          std.write(justValues, '', { format: stdoutFormat });
        } else {
          std.write(justValues[0], '', { format: stdoutFormat });
        }
      } else {
        for (const o of files) {
          const { path, value, ...args } = o;
          std.write(value, path, { overwrite, ...args });
        }
      }
    });
  });
}

export { generate };
