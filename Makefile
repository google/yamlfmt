.EXPORT_ALL_VARIABLES:

.PHONY: build
build:
	go build ./cmd/yamlfmt

.PHONY: test
test:
	go test ./...

.PHONY: test_v
test_v:
	go test -v ./...

YAMLFMT_BIN ?= $(shell pwd)/yamlfmt
.PHONY: export_yamlfmt_bin
export_yamlfmt_bin: export YAMLFMT_BIN = $(YAMLFMT_BIN)

.PHONY: integrationtest_local
integrationtest_local:
	go test -tags=integration_test ./integrationtest/local -update	

.PHONY: integrationtest_local_v
integrationtest_local_v:
	go test -tags=integration_test ./integrationtest/local -update	

.PHONY: integrationtest_local_update
integrationtest_local_update:
	YAMLFMT_BIN="$(YAMLFMT_BIN)" go test -tags=integration_test ./integrationtest/local -update	

.PHONY: install
install:
	go install ./cmd/yamlfmt

.PHONY: install_tools
install_tools:
	go install github.com/google/addlicense@latest

.PHONY: addlicense
addlicense:
	addlicense -c "Google LLC" -l apache .
