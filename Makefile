SHELL=/bin/bash

install: security
	scripts/build/install

security:
	@gosec ./...

clean:
	rm -rf _exe/

static: security
	scripts/build/static

build: security
	scripts/build/dynamic

release-build: security
	scripts/build/release-build

.PHONY: clean static all build
