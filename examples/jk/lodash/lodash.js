import { print } from '@jkcfg/std';
import _ from 'lodash-es';

print(_.defaults({ a: 1 }, { a: 3, b: 2 }));
print(_.partition([1, 2, 3, 4], n => n % 2));
