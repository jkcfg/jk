.PHONY: build-image std-install dep all install test api-reference FORCE

PLUGINS =					\
	plugins/jk-plugin-echo/jk-plugin-echo	\
	$(NULL)

all: jk $(PLUGINS)

VERSION := $(shell git describe --tags)

ifneq ($(RW),yes)
	RO = -mod=readonly
endif

ifeq ($(STATIC),yes)
	A = -a
	TAGS += -tags netgo
	LDFLAGS += -extldflags "-static"
endif

jk: pkg/__std/lib/assets_vfsdata.go FORCE
	GO111MODULE=on go build $(RO) $(A) $(TAGS) -o $@ -ldflags '-X main.Version=$(VERSION) -s -w $(LDFLAGS)'

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

plugins/jk-plugin-echo/jk-plugin-echo: FORCE
	GO111MODULE=on go build $(RO) $(A) $(TAGS) -o $@ -ldflags '-X main.Version=$(VERSION) -s -w $(LDFLAGS)' ./$(@D)

D := $(shell go env GOPATH)/bin
install: all
	mkdir -p $(D)
	cp jk $(D)
	$(foreach p,$(PLUGINS),cp $(p) $(D))

build-image:
	docker build -t quay.io/justkidding/build -f build/Dockerfile build/

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
