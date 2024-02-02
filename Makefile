
.PHONY: none
none:


deps: go.mod
	@echo "+ $@"
	@go mod tidy
ifdef CI
	@git diff --exit-code -- go.mod go.sum || { echo "go.mod/go.sum files were updated after running 'go mod tidy', run this command on your local machine and commit the results." ; exit 1 ; }
endif
	go mod verify
	@touch deps

UNAME_S := $(shell uname -s)
HOST_OS := linux
ifeq ($(UNAME_S),Darwin)
    HOST_OS := darwin
endif

TAG := $(shell ./get-tag)

GOBIN := $(CURDIR)/.gobin
PATH := $(GOBIN):$(PATH)

# Makefile on Mac doesn't pass the updated PATH and GOBIN to the shell
# and so, without the following line, the shell does not end up
# trying commands in $(GOBIN) first.
# See https://stackoverflow.com/a/36226784/3690207
SHELL := env GOBIN=$(GOBIN) PATH=$(PATH) /bin/bash

KUBE_LINTER_BIN := $(GOBIN)/kube-linter

########################################
###### Binaries we depend on ###########
########################################

GOLANGCILINT_BIN := $(GOBIN)/golangci-lint
$(GOLANGCILINT_BIN): deps
	@echo "+ $@"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

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

.PHONY: lint
lint: golangci-lint

####################
## Code generation #
####################

.PHONY: go-generated-srcs
go-generated-srcs: deps
	go generate ./...

.PHONY: generated-docs
generated-docs: go-generated-srcs build
	kube-linter templates list --format markdown > docs/generated/templates.md
	kube-linter checks list --format markdown > docs/generated/checks.md

.PHONY: generated-srcs
generated-srcs: go-generated-srcs generated-docs

#############
## Compile ##
#############


.PHONY: build
build: source-code-archive
	@CGO_ENABLED=0 GOOS=darwin scripts/go-build.sh ./cmd/kube-linter
	@CGO_ENABLED=0 GOOS=linux scripts/go-build.sh ./cmd/kube-linter
	@CGO_ENABLED=0 GOOS=windows scripts/go-build.sh ./cmd/kube-linter
	@mkdir -p "$(GOBIN)"
	@cp "bin/$(HOST_OS)/kube-linter" "$(GOBIN)/kube-linter"
	@chmod u+w "$(GOBIN)/kube-linter"

$(KUBE_LINTER_BIN):
	@$(MAKE) build

.PHONY: image
image: build
	@cp bin/linux/kube-linter image/bin
	@docker build -t "stackrox/kube-linter:$(TAG)" -f image/Dockerfile image/
	@docker build -t "stackrox/kube-linter:$(TAG)-alpine" -f image/Dockerfile_alpine image/

.PHONY: source-code-archive
source-code-archive:
	git archive --prefix="kube-linter-$(TAG)/" HEAD -o "bin/kube-linter-source.tar.gz"
	git archive --prefix="kube-linter-$(TAG)/" HEAD -o "bin/kube-linter-source.zip"


##########
## Test ##
##########

.PHONY: test
test:
	echo $(TAG)
	go test ./...

.PHONY: e2e-test
e2e-test: $(KUBE_LINTER_BIN)
	KUBE_LINTER_BIN="$(KUBE_LINTER_BIN)" go test -tags e2e -count=1 ./e2etests/...

.PHONY: e2e-bats
e2e-bats: $(KUBE_LINTER_BIN)
	@command -v jq &> /dev/null || { echo >&2 'ERROR: jq not installed; See: https://stedolan.github.io/jq/download - Aborting'; exit 1; }
	@command -v diff &> /dev/null || { echo >&2 'ERROR: diff not installed; See: https://www.baeldung.com/linux/diff-command - Aborting'; exit 1; }
	@command -v bats &> /dev/null || { echo >&2 'ERROR: bats not installed; See: https://bats-core.readthedocs.io/en/stable/installation.html - Aborting'; exit 1; }

	KUBE_LINTER_BIN="$(KUBE_LINTER_BIN)" e2etests/bats-tests.sh
	KUBE_LINTER_BIN="$(KUBE_LINTER_BIN)" e2etests/check-bats-tests.sh
