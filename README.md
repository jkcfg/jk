```shell
$ cat nginx.js 
import k from 'kubernetes';
import jk from 'jk';

const container = k.Container('nginx', 'nginx:1.15.4');
const deployment = k.Deployment('nginx', 3, [container]);
jk.write(deployment);
$ jk nginx.js 
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
