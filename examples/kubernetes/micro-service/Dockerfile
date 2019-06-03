FROM alpine:3.8

LABEL maintainer=damien@weave.works \
      org.opencontainers.image.title="myservice" \
      org.opencontainers.image.description="This service has a very useful API" \
      org.opencontainers.image.source="git@github.com:dlespiau/myservice.git" \
      org.opencontainers.image.revision="0.2.0"

WORKDIR /
COPY myservice /
ENTRYPOINT /myservice

EXPOSE  80
