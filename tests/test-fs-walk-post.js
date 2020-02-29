import { walk } from '@jkcfg/std/fs';
import { print } from '@jkcfg/std';

// walk should call post for every time it finishes a directory. to
// check that the hook is called at the right time, include it in the
// transcript.

const post = () => print('post');
const pre = f => print(`pre ${f.name}`) || true;

for (const f of walk('./fs-walk-preorder-files', { pre, post })) {
  if (!f.name.startsWith('.')) print(f.name);
}
