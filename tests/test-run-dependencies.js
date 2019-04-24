import * as std from '@jkcfg/std';
import { foo } from './test-run-dependencies/failure';

std.log(foo);
std.read('test-run-dependencies/svc-myapp.yaml');
