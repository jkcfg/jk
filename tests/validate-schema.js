import { read, log } from '@jkcfg/std';
import { validateWithObject, validateWithFile } from '@jkcfg/std/schema';
import validate from './validate-schema-files/module';

function stringifyResult(result) {
  if (result === 'ok') return result;
  return result.map(err => `${err.path}: ${err.msg}`);
}

export default async function doTest(value) {
  log('Object:');
  const schema = await read('./validate-schema-files/person.json');
  log(stringifyResult(validateWithObject(value, schema)));

  log('File:');
  log(stringifyResult(await validateWithFile(value, 'validate-schema-files/person.json')));

  log('Module:');
  log(stringifyResult(await validate(value)));
}
