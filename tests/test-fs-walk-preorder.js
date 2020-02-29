import { walk } from '@jkcfg/std/fs';
import { print } from '@jkcfg/std';

// The files and directories in fs-walk-preorder-files are designed to
// be in alphabetical order, when traversed in preorder.

// To check the pre hook is called _immediately after_ we've seen a
// directory, this is set for every file, and verified in the
// pre-hook.
let lastname = '';
const pre = (f) => {
  if (lastname !== f.name) {
    throw new Error(`expected to see '${lastname}' in pre-hook, saw '${f.name}'`);
  }
  return true;
};

for (const f of walk('./fs-walk-preorder-files', { pre })) {
  lastname = f.name;
  if (!f.name.startsWith('.')) print(f.name);
}
