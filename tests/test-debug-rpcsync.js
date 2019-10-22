import { echoSync } from '@jkcfg/std/debug';
import { print } from '@jkcfg/std';

print(echoSync());

const arr = new Uint8Array(new ArrayBuffer(3));
arr[0] = 1;
arr[1] = 2;
arr[2] = 3;

print(echoSync(65, arr, 'string', { object: 'object' }));
