import { read, log } from '@jkcfg/std';

read('read-files/deployment-schema.json').then(_ => log('ok'));
