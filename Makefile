SHELL=/bin/bash

GO_LINKER_SYMBOL := "main.version"

default:
	$(MAKE) build

%:
    @:

args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

all: build

clean:
	rm -rf _exe/

static:
	CGO_ENABLED=0 go build -x -ldflags '-w -extldflags "-static"' -o _exe/liquidweb-cli github.com/liquidweb/liquidweb-cli

build:
	go build -ldflags="-s -w" -o _exe/liquidweb-cli github.com/liquidweb/liquidweb-cli

run:
	go run main.go $(call args,)

.PHONY: clean static all build run
