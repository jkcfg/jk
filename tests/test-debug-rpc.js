import { echo } from '@jkcfg/std/debug';
import { log } from '@jkcfg/std';

echo().then(log);

const arr = new Uint8Array(new ArrayBuffer(3));
arr[0] = 1;
arr[1] = 2;
arr[2] = 3;

echo(65, arr, 'string', { object: 'object' }).then(log);
