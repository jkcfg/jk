.PHONY: build-image std-install dep all jk install test

all: jk

VERSION := $(shell git describe)

jk: pkg/__std/lib/assets_vfsdata.go
	GO111MODULE=on go build -o $@ -ldflags "-X main.Version=$(VERSION)"

pkg/__std/lib/assets_vfsdata.go: std/build/std.js
	GO111MODULE=on go generate ./pkg/__std/lib

std/build/std.js: std/*.fbs std/*.js std/package.json
	std/generate.sh
	cd std && npm run build

D := $(shell go env GOPATH)/bin
install: jk
	mkdir -p $(D)
	cp jk $(D)

build-image:
	docker build -t quay.io/justkidding/build -f build/Dockerfile build/

# Pulls the std/node_modules directory
std-install:
	cd std && npm install

# This target install build dependencies
dep: std-install

test:
	./run-tests.sh

clean:
	@rm jk
	@rm -rf .bash_history .cache/ .config/ .npm
