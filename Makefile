PROJECTNAME=$(shell basename "$(PWD)")
GOLANGCI := $(GOPATH)/bin/golangci-lint

.PHONY: help lint test
all: help
help: Makefile
	@echo
	@echo " Choose a make command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

get-lint:
	if [ ! -f ./bin/golangci-lint ]; then \
		wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest; \
	fi;

lint: get-lint
	./bin/golangci-lint run ./... --timeout 5m0s

test:
	./scripts/unit_tests.sh

e2e-setup:
	docker-compose -f ./e2e/evm-optimism/docker-compose-nobuild.yml -f ./e2e/evm-optimism/docker-compose.e2e.yml up --scale verifier=1 -d

e2e-test:
	./scripts/int_tests.sh
