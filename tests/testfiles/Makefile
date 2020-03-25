# Makefile for the images used in automated tests

.PHONY: all

all: foolib.tar barlib.tar

foolib.tar: src/foolib/jk/modules/foolib.js
	tar -C src/foolib --exclude='*~' -c jk | docker import - foolib:v1
	docker save -o foolib.tar foolib:v1

barlib.tar: src/barlib/Dockerfile.base-foo src/barlib/Dockerfile.alpine-foo src/barlib/Dockerfile.alpine-bar
	for image in base-foo alpine-foo alpine-bar; do \
	  docker build -f "src/barlib/Dockerfile.$$image" -t "$$image" ./src/barlib; \
	done
	go run rebase.go --old "alpine-foo" --new "base-foo" --original "alpine-bar" --out "barlib.tar"
