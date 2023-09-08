.EXPORT_ALL_VARIABLES:

.PHONY: build
build:
	go build -o dist/yamlfmt ./cmd/yamlfmt 

.PHONY: test
test:
	go test ./...

.PHONY: test_v
test_v:
	go test -v ./...

YAMLFMT_BIN ?= $(shell pwd)/dist/yamlfmt
.PHONY: integrationtest
integrationtest:
	$(MAKE) build
	go test -tags=integration_test ./integrationtest/command

.PHONY: integrationtest_v
integrationtest_v:
	$(MAKE) build
	go test -v -tags=integration_test ./integrationtest/command

.PHONY: integrationtest_local_update
integrationtest_update:
	$(MAKE) build
	go test -tags=integration_test ./integrationtest/command -update	

.PHONY: install
install:
	go install ./cmd/yamlfmt

.PHONY: install_tools
install_tools:
	go install github.com/google/addlicense@latest

.PHONY: addlicense
addlicense:
	addlicense -c "Google LLC" -l apache .
