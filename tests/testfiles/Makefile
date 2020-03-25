# Makefile for the images used in automated tests

foolib.tar: src/foolib/jk/modules/foolib.js
	tar -C src/foolib --exclude='*~' -c jk | docker import - foolib:v1
	docker save -o foolib.tar foolib:v1
