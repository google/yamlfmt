build:
	go build ./cmd/yamlfmt

install:
	go install ./cmd/yamlfmt

install_tools:
	go install github.com/google/addlicense@latest

addlicense:
	addlicense -c "Google LLC" -l apache .