export TRANSFERWISER_TWLOGINREDIRECT ?= "http://localhost:3000/oauth/callback"

.PHONY: bdd
bdd: godog
	@godog

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


.PHONY: ca_cert
ca_cert:
	openssl genrsa -des3 -out proxy/ca.key 4096
	openssl req -new -x509 -days 365 -key proxy/ca.key -out proxy/ca.crt

.PHONY: server_cert
server_cert:
	openssl genrsa -des3 -out proxy/server.key 1024
	openssl req -new -key proxy/server.key -out proxy/server.csr
	openssl x509 -req -days 365 -in proxy/server.csr -CA proxy/ca.crt -CAkey proxy/ca.key -set_serial 01 -out proxy/server.crt
	openssl rsa -in proxy/server.key -out proxy/temp.key
	rm proxy/server.key
	mv proxy/temp.key proxy/server.key

.PHONY: client_cert
client_cert:
	openssl genrsa -des3 -out proxy/client.key 1024
	openssl req -new -key proxy/client.key -out proxy/client.csr
	openssl x509 -req -days 365 -in proxy/client.csr -CA proxy/ca.crt -CAkey proxy/ca.key -set_serial 01 -out proxy/client.crt
	openssl rsa -in proxy/client.key -out proxy/temp.key
	rm proxy/client.key
	mv proxy/temp.key proxy/client.key
	openssl pkcs12 -export -clcerts -in proxy/client.crt -inkey proxy/client.key -out proxy/client.p12
