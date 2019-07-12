// user/group to run the service as.
const user = service => service.user || 'app';
const group = service => service.group || 'app';
const home = service => service.home || `/home/${user(service)}`;

const baseImage = 'alpine:3.8';

const Dockerfile = service => `FROM ${baseImage}

RUN addgroup -S ${group(service)} \\
    && adduser -D -S -h ${home(service)} -s /sbin/nologin -G ${group(service)} ${user(service)} \\
WORKDIR ${home(service)}
USER ${user(service)}

EXPOSE ${service.port}

COPY ${service.name} ${home(service)}
ENTRYPOINT ${home(service)}/${service.name}
`;

export {
  Dockerfile,
};
