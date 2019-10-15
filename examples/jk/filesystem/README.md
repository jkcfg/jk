# Filesystem examples

This directory has examples of how to explore a filesystem from a `jk`
script. In particular, this demonstrates how to use fs.walk to
enumerate and filter files under a directory.

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
