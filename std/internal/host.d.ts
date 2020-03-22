// This contains definitions for @jkcfg/std/internal/host, which is
// normally supplied by the runtime.

import { ReadOptions } from '../read';
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
export declare function read(path: string, opts?: ReadOptions): Promise<any>;

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
