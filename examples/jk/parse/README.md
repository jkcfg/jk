# Using parse and stringify

This example shows an example of using both of parse and stringify,
from the built-in standard library.

These functions are for parsing strings into values, and serialising
values into strings, respectively. You supply the format you expect
(or want) as the second argument.

## Running the example

The file [`example.yaml`](./example.yaml) is a Kubernetes ConfigMap
definition, with another file embedded in its `data` field.

```bash
$ cat example.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
data:
  config.json: |
    {
      "path": "some/path/somewhere"
    }
```

The script [`add-version.js`](./add-version.js) will transform the
embedded file config.json, to give it a version field:

```bash
$ jk transform --stdout ./add-version.js ./example.yaml
apiVersion: v1
data:
  config.json: '{"path":"some/path/somewhere","version":2}'
kind: ConfigMap
metadata:
  name: config
```

You could also write it to another directory, and compare the files:

```bash
$ jk transform -o /tmp ./add-version.js ./example.yaml
$ diff ./example.yaml /tmp/example.yaml
# ...
```
