import std from 'std';

const raw = std.read('test-read.expected/raw.txt');
raw.then(s => std.write(String.fromCharCode(...s), 'raw.txt', { format: std.Format.Raw }),
         err => std.write(`[ERROR] ${err.toString()}`));

const utf16 = std.read('test-read.expected/utf16.txt', { encoding: std.Encoding.UTF16 });
utf16.then(s => std.write(s, 'utf16.txt', { format: std.Format.Raw }),
           err => std.write(`[ERROR] ${err.toString()}`));
