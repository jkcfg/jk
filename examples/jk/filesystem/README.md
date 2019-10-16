# Filesystem examples

This directory has examples of how to explore a filesystem from a `jk`
script. In particular, this demonstrates how to use fs.walk to
enumerate and filter files under a directory.

## tree.js

This example shows how to print a tree of the files under the input
directory. To print the files under `$CONFIG`:

    jk run -i $CONFIG ./tree.js

This is usually at least a little tricky, since `walk` is turning a
tree structure into a sequence, so we don't have enough information to
tell whether we're going into a directory or returning to a parent
directory.

To be able to recover this information, `walk` has hooks for when it's
ding either of those things, and it always does a [_preorder_
traversal](https://en.wikipedia.org/wiki/Tree_traversal#Pre-order_(NLR)). The
hooks mean you can keep your own state while walking, and the preorder
traversal means you will always see a directory immediately before you
see any files in that directory. (However: the file encountered
immediately after directory A is not necessarily in directory A -- A
could be empty.)

`tree.js` uses a post hook to keep track of where it is in the
directory structure. Whenever it sees a directory, it indents one
level; whenever the post hook is called, it outdents. It uses the pre
hook to filter out dotted directories. (The pre hook is not used to
indent, because it would indent the directory itself.)

## find.js

This example scans the input directory for files passing the
predicates given as parameters. The parameters take the form of
function expressions, with the defaults shown here:

 - `-p match.file="name => name.endsWith('.yaml') || name.endsWith('.yml')"`
 - `-p match.obj="_ => true"`

`match.file` is evaluated on the file name, and `match.obj` is
evaluated on the object loaded from a file.

Here's an example of running find.js to list all the YAML files under
`$CONFIG` that represent Kubernetes deployment resources (NB the
default filename filter is to include YAML files, so no need to supply
that param):

```bash
jk run -i $CONFIG ./find.js -p match.obj="obj => obj.kind == 'Deployment'"
```
