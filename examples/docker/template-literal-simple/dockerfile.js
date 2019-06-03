import * as param from '@jkcfg/std/param';

// input is the the 'service' input parameter.
const input = param.Object('service');

// Our docker images are based on alpine
const baseImage = 'alpine:3.8';

// Dockerfile is a function generating a Dockerfile from a service object.
const Dockerfile = service => `FROM ${baseImage}

EXPOSE ${service.port}

COPY ${service.name} /
ENTRYPOINT /${service.name}`;

// Instruct generate to produce a Dockerfile with the value returned by the
// Dockerfile function.
export default [
  { file: 'Dockerfile', value: Dockerfile(input) },
];
