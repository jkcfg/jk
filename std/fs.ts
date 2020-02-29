/**
 * @module std/fs
 */

import { RPCSync } from './internal/rpc';
import { valueFromUTF8Bytes } from './internal/data';

export class FileInfo {
  name: string;
  path: string;
  isdir: boolean;

  constructor(n: string, p: string, d: boolean) {
    this.name = n;
    this.path = p;
    this.isdir = d;
  }
}

export class Directory {
  name: string;
  path: string;
  files: FileInfo[];

  constructor(n: string, p: string, files: FileInfo[]) {
    this.name = n;
    this.path = p;
    this.files = files;
  }
}

export interface InfoOptions {
  module?: string;
}

export function info(path: string, options: InfoOptions = {}): FileInfo {
  const { module = '' } = options;
  const response = RPCSync('std.fileinfo', path, module);
  const { name: n, path: p, isdir: d } = valueFromUTF8Bytes(response);
  return new FileInfo(n, p, d);
}

export interface DirOptions {
  module?: string;
}

export function dir(path: string, options: DirOptions = {}): Directory {
  const { module = '' } = options;
  const response = RPCSync('std.dir', path, module);
  const { name: n, path: p, files: fs } = valueFromUTF8Bytes(response);
  const infos = [];
  for (const f of fs) {
    const { name: infoname, path: infopath, isdir: infoisdir } = f;
    infos.push(new FileInfo(infoname, infopath, infoisdir));
  }
  return new Directory(n, p, infos);
}

export function join(base: string, name: string): string {
  return `${base}/${name}`;
}

const always = (): boolean => true;
const noop = (): void => {};

export interface WalkOpts {
  pre?: (f: FileInfo) => boolean;
  post?: () => void;
}

/* eslint-disable no-labels */

/** walk is a generator function that yields the files and directories
 * under a given path, in a preorder traversal. "Preorder" means that
 * the traversal is depth-first, and a directory is yielded
 * immediately before the traversal examines its contents.
 *
 * @param path the starting point for the traversal, which should be a
 * directory.
 *
 * @param opts pre- and post-hooks for the walk. The pre-hook is
 * called for each directory after it is yielded; if it returns a
 * falsey value, the directory is not traversed. The post-hook is
 * called after a directory's contents have all been traversed. The
 * starting point does not get treated as part of the traversal, i.e.,
 * it starts with the contents of the directory at `path`.
*/
export function* walk(path: string, opts: WalkOpts = {}): IterableIterator<FileInfo> {
  const { pre = always, post = noop } = opts;
  const top = dir(path);
  // the stack is going to keep lists of files to examine
  const stack: FileInfo[][] = [];
  let next = top.files;

  // eslint-disable-next-line no-restricted-syntax
  runNext: while (next !== undefined) {
    let i = 0;
    for (; i < next.length; i += 1) {
      const f = next[i];
      yield f;
      if (f.isdir && pre(f)) {
        const d = dir(f.path);
        // If we need to recurse into the subdirectory, push the work
        // yet to do here, then process the subdirectory's files.
        if (d.files.length > 0) {
          stack.push(next.slice(i + 1));
          next = d.files;
          continue runNext;
        }
        // If not, we can just continue through the current directory.
        post();
      }
    }
    // If we've exhausted the slice, we're popping a directory
    if (i === next.length && stack.length > 0) post();
    next = stack.pop();
  }
}

/* eslint-enable */
