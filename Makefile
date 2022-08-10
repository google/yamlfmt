install_tools:
	go install github.com/google/addlicense@latest

addlicense:
	addlicense -c "Google LLC" -l apache .