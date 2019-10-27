import * as std from '@jkcfg/std';
import { render } from '@jkcfg/std/render';

render('echo.json', { message: 'success' }).then(r => std.log(r.message));
