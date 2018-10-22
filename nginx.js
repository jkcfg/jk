import k from 'kubernetes';
import jk from 'jk';

const container = k.Container('nginx', 'nginx:1.15.4');
const deployment = k.Deployment('nginx', 3, [container]);
jk.write(deployment);
