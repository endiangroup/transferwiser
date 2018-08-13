
.PHONY: bdd
bdd: godog
	godog

.PHONY: docker_build
docker_build:
	docker build -t endian-group/transferwiser:latest .

# Tools

.PHONY: godog
GODOG_BIN := $(shell command -v godog 2> /dev/null)
godog:
ifndef GODOG_BIN
	@echo "Installing godog..."
	@go get github.com/DATA-DOG/godog/cmd/godog
endif
