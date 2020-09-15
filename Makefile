
.PHONY: none
none:


deps: go.mod
	@echo "+ $@"
	@go mod tidy
	@go mod download
	@go mod verify
	@touch deps


#####################################################################
###### Binaries we depend on ############
#####################################################################

GOBIN := $(CURDIR)/.gobin
PATH := $(GOBIN):$(PATH)
# Makefile on Mac doesn't pass this updated PATH to the shell
# and so without the following line, the shell does not end up
# trying commands in $(GOBIN) first.
# See https://stackoverflow.com/a/36226784/3690207
SHELL := env PATH=$(PATH) /bin/bash

GOLANGCILINT_BIN := $(GOBIN)/golangci-lint
$(GOLANGCILINT_BIN): deps
	@echo "+ $@"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

STATICCHECK_BIN := $(GOBIN)/staticcheck
$(STATICCHECK_BIN): deps
	@echo "+ $@"
	@go install honnef.co/go/tools/cmd/staticcheck

###########
## Lint ##
###########

.PHONY: golangci-lint
golangci-lint: $(GOLANGCILINT_BIN)
ifdef CI
	@echo '+ $@'
	@echo 'The environment indicates we are in CI; running linters in check mode.'
	@echo 'If this fails, run `make lint`.'
	golangci-lint run
else
	golangci-lint run --fix
endif

.PHONY: staticcheck
staticcheck: $(STATICCHECK_BIN)
	staticcheck -checks=all,-ST1000 ./...

.PHONY: lint
lint: golangci-lint staticcheck

#############
## Compile ##
#############

.PHONY: build
build:
	go build -o kube-linter ./cmd/kubelinter
