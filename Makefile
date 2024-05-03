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
	go test -v -tags=integration_test ./integrationtest/command

.PHONY: integrationtest_v
integrationtest_v:
	$(MAKE) build
	go test -v -tags=integration_test ./integrationtest/command

.PHONY: integrationtest_stdout
integrationtest_stdout:
	$(MAKE) build
	go test -v -tags=integration_test ./integrationtest/command -stdout

.PHONY: integrationtest_update
integrationtest_update:
	$(MAKE) build
	go test -tags=integration_test ./integrationtest/command -update

.PHONY: command_test_case
command_test_case:
ifndef TESTNAME
	$(error "TESTNAME undefined")
endif
	mkdir -p integrationtest/command/testdata/$(TESTNAME)/before && \
	mkdir -p integrationtest/command/testdata/$(TESTNAME)/stdout

.PHONY: install
install:
	go install ./cmd/yamlfmt

.PHONY: install_tools
install_tools:
	go install github.com/google/addlicense@latest

.PHONY: addlicense
addlicense:
	addlicense -c "Google LLC" -l apache .

.PHONY: addlicense_check
addlicense_check:
	addlicense -check -c "Google LLC" -l apache .
