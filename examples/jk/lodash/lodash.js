import { log } from '@jkcfg/std';
import _ from 'lodash-es';

log(_.defaults({ a: 1 }, { a: 3, b: 2 }));
log(_.partition([1, 2, 3, 4], n => n % 2));
