import k from 'kubernetes.js';
import std from '@jkcfg/std';

const container = k.Container('nginx', 'nginx:1.15.4');
const deployment = k.Deployment('nginx', 3, [container]);
std.log(deployment, { format: std.Format.YAML });
