import { withModuleRef } from '@jkcfg/std/resource';
import { validateByResource } from '@jkcfg/std/schema';

export default function validate(obj) {
  return withModuleRef(r => validateByResource(obj, 'person.json', r));
}
