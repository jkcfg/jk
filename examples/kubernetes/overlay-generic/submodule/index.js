import { dir, read } from '@jkcfg/std/resource';
import { long } from '@jkcfg/kubernetes/short';

function resources() {
  const ls = dir('.');
  const files = [];
  for (const f of ls.files) {
    if (f.name.endsWith('.yaml')) {
      files.push(read(f.path).then(long));
    }
  }
  return Promise.all(files);
}

export { resources };
