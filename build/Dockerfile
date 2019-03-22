# Build v8
FROM golang:1.12.1 as v8builder
RUN apt-get update && apt-get install -y \
    bzip2 \
    libglib2.0-dev \
    libxml2 \
    xz-utils \
  && rm -rf /var/lib/apt/lists/*
RUN go get -d github.com/jkcfg/v8worker2
RUN cd $GOPATH/src/github.com/jkcfg/v8worker2 \
    && ./build.py

# Build flatbuffer compiler, flatc
FROM golang:1.12.1 as flatc-builder
RUN apt-get update && apt-get install -y \
    cmake \
  && rm -rf /var/lib/apt/lists/*
ENV FLATBUFFERS_VERSION 1.10.0
RUN curl -fsSLO --compressed "https://github.com/google/flatbuffers/archive/v${FLATBUFFERS_VERSION}.tar.gz" \
    && tar -xf v${FLATBUFFERS_VERSION}.tar.gz \
    && rm v${FLATBUFFERS_VERSION}.tar.gz \
    && cd flatbuffers-${FLATBUFFERS_VERSION} \
    && cmake -G "Unix Makefiles" \
    && make \
    && cp flatc /usr/local/bin \
    && cd .. \
    && rm -rf flatbuffers-${FLATBUFFERS_VERSION}

# Build github-release
FROM golang:1.12.1 as github-release-builder
RUN go get github.com/aktau/github-release \
  && cp `go env GOPATH`/bin/github-release /usr/local/bin \
  && rm -rf `go env GOPATH`/src/github.com/aktau/github-release

FROM golang:1.12.1 as fetcher
RUN apt-get update && apt-get install -y \
    xz-utils \
  && rm -rf /var/lib/apt/lists/*

# Fetch node and npm
ENV NODE_VERSION 8.12.0
ENV ARCH x64
RUN curl -fsSLO --compressed "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-$ARCH.tar.xz" \
    && tar -xJf "node-v$NODE_VERSION-linux-$ARCH.tar.xz" -C /usr/local --strip-components=1 --no-same-owner

# Fetch gometalinter
ENV GOMETALINTER_VERSION 2.0.11
ENV ARCH amd64
RUN curl -fsSLO --compressed "https://github.com/alecthomas/gometalinter/releases/download/v$GOMETALINTER_VERSION/gometalinter-$GOMETALINTER_VERSION-linux-$ARCH.tar.gz" \
    && tar -xf "gometalinter-$GOMETALINTER_VERSION-linux-$ARCH.tar.gz" -C /usr/local/bin --strip-components=1 --no-same-owner

# Fetch gosu
ENV GOSU_VERSION 1.11
ENV ARCH amd64
RUN curl -fsSLo /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$ARCH" \
    && chmod +x /usr/local/bin/gosu

# Our final build image
FROM golang:1.12.1
COPY --from=v8builder /go/src/github.com/jkcfg/v8worker2/v8.pc /usr/local/lib/pkgconfig/
RUN sed -i \
     -e 's#Cflags: .*#Cflags: -I/usr/local/include/v8#' \
     -e 's#Libs: .*#Libs: /usr/local/lib/v8/libv8_monolith.a#' \
     /usr/local/lib/pkgconfig/v8.pc
COPY --from=v8builder /go/src/github.com/jkcfg/v8worker2/v8/include /usr/local/include/v8/
COPY --from=v8builder /go/src/github.com/jkcfg/v8worker2/out/v8build/obj/libv8_monolith.a /usr/local/lib/v8/
COPY --from=flatc-builder /usr/local/bin/ /usr/local/bin/
COPY --from=github-release-builder /usr/local/bin/ /usr/local/bin/
COPY --from=fetcher /usr/local/bin/ /usr/local/bin/
COPY --from=fetcher /usr/local/lib/node_modules/ /usr/local/lib/node_modules/

ENV SRC_PATH /go/src/github.com/jkcfg/jk
WORKDIR $SRC_PATH
COPY entrypoint.sh /usr/local/bin
ENTRYPOINT ["entrypoint.sh"]
CMD ["/bin/bash"]
