const myService = {
  name: 'my-service',
  port: 80,
};

const baseImage = 'alpine:3.8';

const Dockerfile = service => `FROM ${baseImage}

EXPOSE ${service.port}

COPY ${service.name} /
ENTRYPOINT /${service.name}`;

export default [
  { file: 'Dockerfile', value: Dockerfile(myService) },
];
