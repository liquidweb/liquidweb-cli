SHELL=/bin/bash

default:
	$(MAKE) install

clean:
	rm -rf _exe/

static:
	scripts/build/static

build:
	scripts/build/dynamic

install:
	scripts/build/install

release-build:
	scripts/build/release-build

.PHONY: clean static all build
