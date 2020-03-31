// This contains definitions for @jkcfg/std/internal/host, which is
// normally supplied by the runtime.

import { ReadOptions } from '../read';
import { WriteOptions } from '../write';
import {
  InfoOptions,
  DirOptions,
  FileInfo,
  Directory,
} from '../fs';

/**
 * read shall read the contents of the file at the path given,
 * relative to the importing module.
 */
// NB this is not the same signature as std.read, because it is not
// intended for reading from stdin.
export declare function read(path: string, opts?: ReadOptions): Promise<any>;

/**
 * write shall write the value to the host filesystem at the path
 * given.
 */
// NB this is not quite the same signature as std.write, because it is
// not intended for writing to stdout.
export declare function write(value: any, path: string, opts?: WriteOptions): void;

/**
 * info shall give the file info at the path given, relative to the
 * importing module.
 */
export declare function info(path: string, options: InfoOptions): FileInfo;

/**
 * dir shall give the directory info at the path given, relative to
 * the importing module.
 */
export declare function dir(path: string, options: DirOptions): Directory;

/**
 * @ignore
 */
export declare function withModuleRef(fn: (ref: string) => any): any;
