import { dir, info } from '@jkcfg/std/resource';

export default Promise.resolve({ dir: dir('dir'), info: info('dir/info') });
