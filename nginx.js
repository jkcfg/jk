import k from 'kubernetes.js';
import std from 'std';

const container = k.Container('nginx', 'nginx:1.15.4');
const deployment = k.Deployment('nginx', 3, [container]);
std.log(deployment, std.Format.YAML);
