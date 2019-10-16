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

const always = (_: FileInfo): boolean => true;
const noop = (): void => {};

interface WalkOpts {
  pre?: (f: FileInfo) => boolean;
  post?: () => void;
}

export function* walk(path: string, opts: WalkOpts = {}): IterableIterator<FileInfo> {
  const { pre = always, post = noop } = opts;
  const top = dir(path);
  // the stack is going to keep lists of files to examine
  const stack: FileInfo[][] = [];
  let next = top.files;
  while (next !== undefined) {
    let i = 0;
    for (; i < next.length; i += 1) {
      const f = next[i];
      if (f.isdir && pre(f)) {
        const d = dir(f.path);
        // If we need to recurse into the subdirectory, push the work
        // yet to do here, then the subdirectory's files. If not, we
        // can just continue as before.
        if (d.files.length > 0) {
          if (i < next.length - 1) stack.push(next.slice(i + 1));
          stack.push(d.files);
          yield f;
          break;
        } else {
          yield f;
          post();
          continue;
        }
      }
      yield f;
    }
    // If we've exhausted the slice, we're popping a directory
    if (i === next.length) post();
    next = stack.pop();
  }
}
