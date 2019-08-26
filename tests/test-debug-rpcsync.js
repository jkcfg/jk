import { echoSync } from '@jkcfg/std/debug';
import { log } from '@jkcfg/std';

log(echoSync());

const arr = new Uint8Array(new ArrayBuffer(3));
arr[0] = 1;
arr[1] = 2;
arr[2] = 3;

log(echoSync(65, arr, 'string', { object: 'object' }));
