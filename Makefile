export TRANSFERWISER_TWAPITOKEN ?= abcd123
export TRANSFERWISER_TWPROFILE ?= 123
export TRANSFERWISER_TWAPITOKEN ?= 12345678-1234-1234-1234-123456789012

.PHONY: test
test: 
	go test ./...

.PHONY: docker_build
docker_build:
	docker build -t endian-group/transferwiser:latest .

.PHONY: ca_cert
ca_cert:
	openssl genrsa -des3 -out certs/ca.key 4096
	openssl req -x509 -new -nodes -key certs/ca.key -sha256 -days 1024 -out certs/ca.crt

.PHONY: server_cert
server_cert:
	openssl genrsa -des3 -out certs/server.key 1024
	openssl req -new -key certs/server.key -out certs/server.csr
	openssl x509 -req -days 365 -in certs/server.csr -CA certs/ca.crt -CAkey certs/ca.key -set_serial 01 -out certs/server.crt
	openssl rsa -in certs/server.key -out certs/temp.key
	rm certs/server.key
	mv certs/temp.key certs/server.key

.PHONY: client_cert
client_cert:
	openssl genrsa -des3 -out certs/client.key 1024
	openssl req -new -key certs/client.key -out certs/client.csr
	openssl x509 -req -days 365 -in certs/client.csr -CA certs/ca.crt -CAkey certs/ca.key -set_serial 01 -out certs/client.crt
	openssl rsa -in certs/client.key -out certs/temp.key
	rm certs/client.key
	mv certs/temp.key certs/client.key
	openssl pkcs12 -export -clcerts -in certs/client.crt -inkey certs/client.key -out certs/client.p12
	openssl pkcs12 -in certs/client.p12 -out certs/client.pem

.PHONY: dep
DEP_BIN := $(shell command -v dep 2> /dev/null)
dep:
ifndef DEP_BIN
	@echo "Installing dep..."
	@go get github.com/golang/dep/cmd/dep
endif
