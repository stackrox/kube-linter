
.PHONY: all
all:


.PHONY: deps
deps:
	@echo "+ $@"
	@go mod tidy
	@go mod download
	@go mod verify


#####################################################################
###### Binaries we depend on ############
#####################################################################

GOLANGCILINT_BIN := $(GOPATH)/bin/golangci-lint
$(GOLANGCILINT_BIN): deps
	@echo "+ $@"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

STATICCHECK_BIN := $(GOPATH)/bin/staticcheck
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
	staticcheck -checks=all ./...

.PHONY: lint
lint: golangci-lint staticcheck
