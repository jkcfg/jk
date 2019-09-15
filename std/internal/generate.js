import * as std from '@jkcfg/std';
import * as param from '@jkcfg/std/param';
import generateDefinition from '%s';

const inputParams = {
  stdout: param.Boolean('jk.generate.stdout', false),
};

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

function error(msg) {
  std.log(`error: ${msg}`);
}

function help() {
  std.log(helpMsg);
}

/**
 * Calculate the modulus of two numbers
 * @param {number} x
 * @param {number} y
 * @returns {number} res
 * @private
 */
function mod(x, y) {
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
const nth = (n) => {
  const s = ['th', 'st', 'nd', 'rd'];

  const v = mod(n, 100);
  return n + (s[mod(v - 20, 10)] || s[v] || s[0]);
};

function extension(path) {
  return path.split('.').pop();
}

function formatFromPath(path) {
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

const isString = s => typeof s === 'string' || s instanceof String;

// Compute the output format of a value.
function valueFormat(o) {
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

function formatSummary(value) {
  const formats = Array(Object.keys(std.Format).length).fill(0);

  value.forEach((e) => {
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

const formatName = f => formatNames[f];

function usedFormats(summary) {
  const augmented = summary.map((n, i) => ({ format: formatName(i), n }));
  return augmented.reduce((formats, desc) => {
    if (desc.n > 0) {
      formats.push(desc.format);
    }
    return formats;
  }, []);
}

function validate(value, params) {
  /* we have an array */
  if (!Array.isArray(value)) {
    error('default value is not an array');
    return { valid: false, showHelp: true };
  }

  /* an array with each element a { path, value } object */
  let valid = true;
  value.forEach((e, i) => {
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

  if (valid === false) {
    return { valid, showHelp: true };
  }

  /* when outputting to stdout, ensure that: */
  let stdoutFormat;
  if (params.stdout === true) {
    /* there's a single output format defined */
    const summary = formatSummary(value);
    const formats = usedFormats(summary);
    if (formats.length > 1) {
      error(`stdout output requires using a single format but got: ${formats.join(',')}`);
      return { valid: false, showHelp: false };
    }

    /*
     * If we have more than one file to generate, make sure it's either JSON or
     * YAML so we can output a stream of documents.
     */
    if (value.length > 1 && formats[0] !== 'JSON' && formats[0] !== 'YAML') {
      error(`stdout output for multiple files requires either JSON or YAML format but got: ${formats[0]}`);
      return { valid: false, showHelp: false };
    }

    if (value.length > 1) {
      if (formats[0] === 'JSON') {
        stdoutFormat = std.Format.JSONStream;
      } else if (formats[0] === 'YAML') {
        stdoutFormat = std.Format.YAMLStream;
      }
    } else {
      stdoutFormat = valueFormat(value[0]);
    }
  }

  return { valid: true, stdoutFormat, showHelp: false };
}


function generate(defaultExport, params) {
  /*
   * The default export can be:
   *  1. an array of { path, value } objects,
   *  2. a promise to such an array,
   *  3. a function evaluating to either 1. or 2.
   */
  let definition = defaultExport;
  if (typeof definition === 'function') {
    definition = definition();
  }

  Promise.resolve(definition).then((files) => {
    /* values can be promises as well */
    const values = files.map(f => f.value);
    Promise.all(values).then(resolved => {
      resolved.map((v, i) => files[i].value = v);

      const { valid, stdoutFormat, showHelp } = validate(files, params);
      if (showHelp) {
        help();
      }
      if (!valid) {
        throw new Error('jk-internal-skip: validation failed');
      }

      if (params.stdout) {
        if (files.length > 1) {
          const values = files.map(f => f.value);
          std.write(values, '', { format: stdoutFormat });
        } else {
          std.write(files[0].value, '', { format: stdoutFormat });
        }
      } else {
        for (const o of files) {
          const { path, value, ...args } = o;
          std.write(value, path, args);
        }
      }
    })
  });
}

generate(generateDefinition, inputParams);
