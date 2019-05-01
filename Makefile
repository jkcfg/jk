.PHONY: build-image std-install dep all install test FORCE

all: jk

VERSION := $(shell git describe --tags)

jk: pkg/__std/lib/assets_vfsdata.go FORCE
	GO111MODULE=on go build -o $@ -ldflags "-X main.Version=$(VERSION)"

pkg/__std/lib/assets_vfsdata.go: std/__std_generated.ts std/dist/index.js
	GO111MODULE=on go generate ./pkg/__std/lib

std/__std_generated.ts: std/*.fbs std/package.json std/generate.sh
	std/generate.sh

std/dist/index.js: std/*.js std/*.ts
	cd std && npm run build

module = @jkcfg/std
module: $(module)/package.json
$(module)/package.json: std/*.js std/*.ts std/__std_generated.ts std/package.json
	cd std && npx tsc --outDir ../$(module)
	cd std && npx tsc --declaration --emitDeclarationOnly --allowJs false --outdir ../$(module) || true
	cp README.md LICENSE std/package.json std/flatbuffers.d.ts $(module)

D := $(shell go env GOPATH)/bin
install: jk
	mkdir -p $(D)
	cp jk $(D)

build-image:
	docker build -t quay.io/justkidding/build -f build/Dockerfile build/

# Pulls the std/node_modules directory
std-install:
	cd std && npm ci

# This target install build dependencies
dep: std-install

test: module
	./run-tests.sh

clean-tests:
	@rm -rf tests/*.got

clean: clean-tests
	@rm -f jk
	@rm -rf .bash_history .cache/ .config/ .npm
	@rm -rf std/dist std/__std_generated.js std/__std_generated.ts
	@rm -rf @jkcfg
