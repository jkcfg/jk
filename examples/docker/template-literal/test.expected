FROM alpine:3.8

RUN addgroup -S app \
    && adduser -D -S -h /home/app -s /sbin/nologin -G app app \
WORKDIR /home/app
USER app

EXPOSE 80

COPY my-service /home/app
ENTRYPOINT /home/app/my-service
