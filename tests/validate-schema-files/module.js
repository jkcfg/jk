import { withModuleRef } from '@jkcfg/std/resource';
import { validateWithResource } from '@jkcfg/std/schema';

export default function validate(obj) {
  return withModuleRef(r => validateWithResource(obj, 'person.json', r));
}
