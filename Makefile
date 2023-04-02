.PHONY: build
build:
	go build ./cmd/yamlfmt

.PHONY: test
test:
	go test ./...

.PHONY: test_v
test_v:
	go test -v ./...

.PHONY: install
install:
	go install ./cmd/yamlfmt

.PHONY: install_tools
install_tools:
	go install github.com/google/addlicense@latest

.PHONY: addlicense
addlicense:
	addlicense -c "Google LLC" -l apache .
