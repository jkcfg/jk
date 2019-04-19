import std from '@jkcfg/std';

std.log('success');
std.write('success', '');
std.write('success', '', { format: std.Format.JSON, indent: 2, overwrite: false });
std.read('success.json').then(json => std.log(json.message));
