.PHONY: build-image std-install dep all install test api-reference FORCE

all: jk

VERSION := $(shell git describe --tags)

ifneq ($(RW),yes)
	RO = -mod=readonly
endif

jk: pkg/__std/lib/assets_vfsdata.go FORCE
ifeq ($(STATIC),yes)
	GO111MODULE=on go build $(RO) -a -tags netgo -o $@ -ldflags '-X main.Version=$(VERSION) -s -w -extldflags "-static"'
else
	GO111MODULE=on go build $(RO) -o $@ -ldflags "-X main.Version=$(VERSION) -s -w"
endif

pkg/__std/lib/assets_vfsdata.go: std/internal/__std_generated.ts std/dist/index.js
	GO111MODULE=on go generate $(RO) ./pkg/__std/lib

std/internal/__std_generated.ts: std/internal/*.fbs std/package.json std/generate.sh
	std/generate.sh

std_sources = std/*.js std/*.ts std/internal/*.ts std/internal/*.js std/cmd/*.ts std/cmd/*.js

std/dist/index.js: $(std_sources)
	rm -rf ./std/dist
	mkdir -p std/dist
	cd std && npm run build

module = @jkcfg/std
module: $(module)/package.json
$(module)/package.json: $(std_sources) std/internal/__std_generated.ts std/package.json
	cd std && npx tsc --outDir ../$(module)
	cd std && npx tsc --declaration --emitDeclarationOnly --allowJs false --outdir ../$(module) || true
	cp README.md LICENSE std/package.json std/internal/flatbuffers.d.ts $(module)

D := $(shell go env GOPATH)/bin
install: jk
	mkdir -p $(D)
	cp jk $(D)

build-image:
	docker build -t jkcfg/build -f build/Dockerfile build/

# Pulls the std/node_modules directory
std-install:
	cd std && npm ci

# Clone the theme directory
typedoc-theme-install:
	rm -rf std/typedoc-theme && git clone https://github.com/jkcfg/typedoc-theme.git std/typedoc-theme

# This target installs build dependencies
dep: std-install typedoc-theme-install

test: module
	./run-tests.sh

api-reference: $(std_sources)
	cd std && npm run doc

clean-tests:
	@rm -rf tests/*.got

clean: clean-tests
	@rm -f jk
	@rm -rf .bash_history .cache/ .config/ .npm
	@rm -rf std/dist std/internal/__std_generated.js std/internal/__std_generated.ts
	@rm -rf @jkcfg

dep-clean: clean
	@rm -rf std/node_modules
	@rm -rf std/typedoc-theme
