import { walk } from '@jkcfg/std/fs';
import { print } from '@jkcfg/std';

// walk should call post for every time it finishes a directory.

const nested = [];
for (const f of walk('./fs-walk-preorder-files', { pre: v => nested.push(v.name), post: () => nested.pop() })) {
  if (!f.name.startsWith('.')) print(f.name);
}

if (nested.length > 0) throw new Error(`did not pop as many directories as pushed, left with ${nested.join(', ')}`);
