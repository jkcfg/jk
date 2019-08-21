import * as std from '@jkcfg/std';
import * as fs from '@jkcfg/std/fs';

const info = fs.info('testfs/foo.txt');
std.write(info, 'fileinfo.json');

const dir = fs.dir('testfs');
std.write(dir, 'dir.json');

const dir2 = fs.dir('testfs/bar');
const bartxt = dir2.files[0];
std.write(bartxt, 'barinfo.json');
