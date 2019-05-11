import * as std from '@jkcfg/std';
import * as param from '@jkcfg/std/param';
import generateDefinition from '%s';

const inputParams = {
  format: param.Number('format', std.Format.FromExtension),
  stdout: param.Boolean('stdout', false),
};

const helpMsg = `
To use generate, export a default value with the list of files to generate:

  export default [
    { file: 'file1.yaml', content: value1 },
    { file: 'file2.yaml', content: [v0, v1, v2], format: std.Format.YAMLStream },
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

function validate(value) {
  /* we have an array */
  if (!Array.isArray(value)) {
    error('default value is not an array');
    return false;
  }

  /* an array with each element a { file, content } object */
  let valid = true;
  value.forEach((e, i) => {
    ['file', 'content'].forEach((prop) => {
      if (!Object.prototype.hasOwnProperty.call(e, prop)) {
        error(`${nth(i + 1)} element does not have a '${prop}' property`);
        valid = false;
      }
    });
  });

  if (valid === false) {
    return false;
  }

  return true;
}


function generate(definition) {
  Promise.resolve(definition).then((files) => {
    if (!validate(files)) {
      help();
      throw new Error('jk-internal-skip: validation failed');
    }

    for (const o of files) {
      const { file, content, ...args } = o;
      std.write(content, file, args);
    }
  });
}

generate(generateDefinition, inputParams);
