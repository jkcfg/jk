import std from 'std';

std.print('success');
std.log('success');
std.write('success', '');
std.write('success', '', { format: std.Format.JSON, indent: 2, override: false });
std.read('success.json').then(json => std.log(json.message));
