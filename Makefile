deps: go.mod go.sum tool-imports/go.sum tool-imports/go.mod
	@touch deps

%.sum: %.mod
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

GOBIN := $(CURDIR)/.gobin
DIST := $(CURDIR)/dist
PATH := $(DIST):$(GOBIN):$(PATH)

KUBE_LINTER_BIN := $(GOBIN)/kube-linter

COVFILES := $(shell mktemp -d)

########################################
###### Binaries we depend on ###########
########################################

GOLANGCILINT_BIN := $(GOBIN)/golangci-lint
$(GOLANGCILINT_BIN): deps
	@echo "+ $@"
	cd tool-imports; \
	GOBIN=$(GOBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint

GORELEASER_BIN := $(GOBIN)/goreleaser
$(GORELEASER_BIN): deps
	@echo "+ $@"
	cd tool-imports; \
	GOBIN=$(GOBIN) go install github.com/goreleaser/goreleaser

###########
## Lint ##
###########

.PHONY: golangci-lint
golangci-lint: $(GOLANGCILINT_BIN)
ifdef CI
	@echo '+ $@'
	@echo 'The environment indicates we are in CI; running linters in check mode.'
	@echo 'If this fails, run `make lint`.'
	$(GOLANGCILINT_BIN) run
else
	$(GOLANGCILINT_BIN) run --fix
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
generated-docs: go-generated-srcs $(KUBE_LINTER_BIN)
	$(KUBE_LINTER_BIN) templates list --format markdown > docs/generated/templates.md
	$(KUBE_LINTER_BIN) checks list --format markdown > docs/generated/checks.md

.PHONY: generated-srcs
generated-srcs: go-generated-srcs generated-docs

#############
## Compile ##
#############


.PHONY: build
build: $(KUBE_LINTER_BIN)

$(KUBE_LINTER_BIN): $(GORELEASER_BIN) $(shell find . -type f -name '*.go')
	$(GORELEASER_BIN) build --snapshot --clean
	mkdir -p $(GOBIN)
	@cp "$(DIST)/kube-linter_$(HOST_OS)_amd64_v1/kube-linter" "$(GOBIN)/kube-linter"
	@chmod u+w "$(GOBIN)/kube-linter"

##########
## Test ##
##########

.PHONY: test
test:
	go test ./... -race -covermode=atomic -coverprofile=coverage.out

.PHONY: e2e-test
e2e-test: $(KUBE_LINTER_BIN)
	KUBE_LINTER_BIN="$(KUBE_LINTER_BIN)" go test -tags e2e -count=1 ./e2etests/...

.PHONY: e2e-bats
e2e-bats: $(KUBE_LINTER_BIN)
	@command -v jq &> /dev/null || { echo >&2 'ERROR: jq not installed; See: https://stedolan.github.io/jq/download - Aborting'; exit 1; }
	@command -v diff &> /dev/null || { echo >&2 'ERROR: diff not installed; See: https://www.baeldung.com/linux/diff-command - Aborting'; exit 1; }
	@command -v bats &> /dev/null || { echo >&2 'ERROR: bats not installed; See: https://bats-core.readthedocs.io/en/stable/installation.html - Aborting'; exit 1; }

	GOCOVERDIR=$(COVFILES) KUBE_LINTER_BIN="$(KUBE_LINTER_BIN)" e2etests/bats-tests.sh
	GOCOVERDIR=$(COVFILES) KUBE_LINTER_BIN="$(KUBE_LINTER_BIN)" e2etests/check-bats-tests.sh
	go tool covdata textfmt -i=$(COVFILES) -o coverage.out
