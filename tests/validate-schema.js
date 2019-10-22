import { read, print } from '@jkcfg/std';
import { validateWithObject, validateWithFile } from '@jkcfg/std/schema';
import validate from './validate-schema-files/module';

function stringifyResult(result) {
  if (result === 'ok') return result;
  return result.map(err => `${err.path}: ${err.msg}`);
}

export default async function doTest(value) {
  print('Object:');
  const schema = await read('./validate-schema-files/person.json');
  print(stringifyResult(validateWithObject(value, schema)));

  print('File:');
  print(stringifyResult(await validateWithFile(value, 'validate-schema-files/person.json')));

  print('Module:');
  print(stringifyResult(await validate(value)));
}
