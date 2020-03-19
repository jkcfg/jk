# Injecting parameters into a template

This example shows how to get parameter values from the command-line
and create a config file, by injecting them into an object then
printing the object.

To run the example (from this directory):

```bash
$ jk run ./template.js -p port=3456 -p name=foo-svc
```

You can also **supply parameters from a file**:

```bash
$ cat > params.yaml <<EOF
port: 6543
name: bar-svc
selector:
  foo: bar
EOF
$ jk run ./template.js -f params.yaml
```

.. or get them from a file, then **override them individually**:

```
$ jk run ./template.js -f params.yaml -p name=foo-svc
```

The parameter values are merged, in the order they are given. So if
you supply `-p name=` before the file, the value in the file is used:

```
$ jk run ./template.js -p name=foo-svc -f params.yaml
```

**Objects** values can be given by **using a nested path**:

```
# this implies select = { app: "helloworld" }
$ jk run ./template.js -p selector.app=helloworld
```

.. or, as in the file example, by giving an object in a parameters
file.

## How this works

The bulk of the script is simply constructing an object literal to be
printed out:

```
const obj = {
  // ...
}

print(obj, { format: Format.YAML });
```

The param module is used to inject the parameters given on the command
line, where they go in the structure:

```
import * as param from '@jkcfg/std/param';

const obj = {
  // ...
  metadata: {
    name: param.String('name', 'service'),
  },
  // ...
};
```

The param methods are given defaults, so there will be a value even if
the parameter named has not been supplied.

Because values are merged, starting with the defaults, it is easy to
not get what you intend when using `param.Object` with a default --
any value supplied will be merged into the default value, so you can
end up with two fields set, one from the default and one from the
parameter value.

This script avoids that by using `false` as the default value, and
substituting another value if no parameter overrides it:

```
const selector = param.Object('selector', false);

const obj = {
  // ... spec: {
    selector: selector || { app: 'app' },
  },
}
```
