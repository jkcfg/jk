[![Build Status](https://travis-ci.org/dlespiau/jk.svg?branch=master)](https://travis-ci.org/dlespiau/jk)

```shell
$ cat nginx.js
import k from 'kubernetes.js';
import std from 'std';

const container = k.Container('nginx', 'nginx:1.15.4');
const deployment = k.Deployment('nginx', 3, [container]);
std.log(deployment, { format: std.Format.YAML });

$ jk run nginx.js
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    containers:
    - image: nginx:1.15.4
      name: nginx
    metadata:
      labels:
        app: nginx
```
