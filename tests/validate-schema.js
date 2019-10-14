import { read, log } from '@jkcfg/std';
import { validateBySchema, validateByFile } from '@jkcfg/std/schema';
import validate from './validate-schema-files/module';

export default async function doTest(value) {
  log('Object:');
  const schema = await read('./validate-schema-files/person.json');
  log(validateBySchema(value, schema));

  log('File:');
  log(await validateByFile(value, 'validate-schema-files/person.json'));

  log('Module:');
  log(await validate(value));
}
