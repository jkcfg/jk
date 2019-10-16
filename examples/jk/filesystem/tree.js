// Prints out a tree picture of the filesystem structure under the
// input directory.

import { log } from '@jkcfg/std';
import { walk } from '@jkcfg/std/fs';

function tree() {
  const stack = [''];
  const pre = dir => !dir.name.startsWith('.');
  const post = () => stack.shift();

  for (const file of walk('.', { pre, post })) {
    log(`${stack[0]}${file.name}`);
    if (file.isdir) {
      stack.unshift(`${stack[0]}  `);
    }
  }
}

tree();
