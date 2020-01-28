SHELL=/bin/bash

GO_LINKER_SYMBOL := "main.version"

default:
	$(MAKE) install

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

install:
	go install
	@echo ""
	@echo "liquidweb-cli has been installed, and it should now be in your PATH."
	@echo ""
	@echo "Executables are installed in the directory named by the GOBIN environment"
	@echo "variable, which defaults to GOPATH/bin or HOME/go/bin if the GOPATH"
	@echo "environment variable is not set. Executables in GOROOT"
	@echo "are installed in GOROOT/bin or GOTOOLDIR instead of GOBIN."
	@echo ""

run:
	go run main.go $(call args,)

release-build:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: clean static all build run
