// This example shows how to use the filesystem part of jk's standard
// library to find YAML files under a directory.
import { read, log } from '@jkcfg/std';
import * as param from '@jkcfg/std/param';
import { walk } from '@jkcfg/std/fs';

/* eslint-disable no-new-func, no-await-in-loop */

const defaultPredicate = '_ => true';
const makePredicate = new Function(`return (${param.String('match.obj', defaultPredicate)});`);

const defaultFilter = 'name => name.endsWith(".yaml") || name.endsWith(".yml")';
const makeFilter = new Function(`return (${param.String('match.file', defaultFilter)});`);

async function find() {
  const top = param.String('path', '.');
  const pred = makePredicate();
  const filter = makeFilter();

  for (const f of walk(top)) {
    if (!f.isdir && filter(f.name)) {
      const obj = await read(f.path);
      if (pred(obj)) {
        log(f.path);
      }
    }
  }
}

find();
