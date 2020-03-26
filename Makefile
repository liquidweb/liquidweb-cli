SHELL=/bin/bash

install: security
	scripts/build/install

security:
	go get github.com/securego/gosec/cmd/gosec
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
