import * as std from '@jkcfg/std';

std.log('success');
std.write('success', std.stdout);
std.write('success', std.stdout, { format: std.Format.JSON, indent: 2, overwrite: false });
std.read('success.json').then(json => std.log(json.message));
