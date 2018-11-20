.PHONY: std-install dep all jk install test

all: jk

jk: pkg/__std/lib/assets_vfsdata.go
	GO111MODULE=on go build -o $@

pkg/__std/lib/assets_vfsdata.go: std/build/std.js
	GO111MODULE=on go generate ./pkg/__std/lib

std/build/std.js: std/*.fbs std/*.js std/package.json
	std/generate.sh
	cd std && npm run build

install: jk
	cp jk `go env GOPATH`/bin

# Pulls the std/node_modules directory
std-install:
	cd std && npm install

# This target install build dependencies
dep: std-install

test:
	./run-tests.sh
