import * as std from '@jkcfg/std';
import { foo } from './test-run-dependencies/failure';

std.print(foo);
std.read('test-run-dependencies/svc-myapp.yaml');
