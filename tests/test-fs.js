import std from '@jkcfg/std';

const info = std.info('testfs/foo.txt');
std.write(info, 'fileinfo.json');

const dir = std.dir('testfs');
std.write(dir, 'dir.json');

const dir2 = std.dir('testfs/bar');
const bartxt = dir2.files[0];
std.write(bartxt, 'barinfo.json');
