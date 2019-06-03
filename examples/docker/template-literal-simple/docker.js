const baseImage = 'alpine:3.8';

const Dockerfile = service => `FROM ${baseImage}

EXPOSE ${service.port}
WORKDIR /

COPY ${service.name} /
ENTRYPOINT /${service.name}`;

export {
  Dockerfile,
};
