export TRANSFERWISER_TWAPITOKEN ?= abcd123
export TRANSFERWISER_TWPROFILE ?= 123
export TRANSFERWISER_TWAPITOKEN ?= 12345678-1234-1234-1234-123456789012
export TRANSFERWISER_CACERT ?= -----BEGIN CERTIFICATE-----MIIEljCCAn4CCQCvmWQhMeNAHDANBgkqhkiG9w0BAQsFADANMQswCQYDVQQGEwJVSzAeFw0xODA5MDUwMzE4MzBaFw0yMTA2MjUwMzE4MzBaMA0xCzAJBgNVBAYTAlVLMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAra+ul375AyFgMo6yMtQpoElDe3Y7+xxmFYIQX1vAH0QpJobNneBQCZ938Hf2NMuENP9znJwnlW6Nv65NUWozeH2mLRIpYBYZOxFm8DIQExo+GW5FQrKIM8q+Ic3ZEGm+bHIIN/ORIHfEU79TUBolpmLBu3wse1RLag2DUCFw4IgieE05PSbXXeiaWW7gqPyASvF4gwy6nzdmQ2HAI4OJgf9+s2WAyQR8nPIOx2AatFYu8ChzbcNsp0d6Z+9CRXlmBrtm0Pr4/HQC3i1BsGyqoKJHewZd4oFSHRTvUonA19EHivdk3fIIlXKA/lUS4/+a38iOOyqh06zDgr2OJiFMQTKyJsHJmrfIRjYVARbPHe5iPgrBF4PzXnFh/Za8ASYpaqMQw7P4S+DpG/4MWg96qDGWHxUqjOeq5hIoDqbAUw5QkVieDDTGmjWO6Xr9ldkcJMl88XXhu0J+KBISmY0apLLF6aYKHE1LrdWOthzgwrdnabacHzieTfC0kXTaxf3sDEgyK/etlpd7RQU0r1B5Qr0OTLDKF2t/KglU8IvWei7WwFn/z/9VdlXvE24jL8imYTSKHSKsMe5uKET1gCKZh1xeZ4B8xxOo9w5iYeHl5wHZaRG2+XnC9sIqwOspLcHMa8vvsy1Go8kPgpe447fZOj8t0q5K4DllSXFf0BQ9rLMCAwEAATANBgkqhkiG9w0BAQsFAAOCAgEAqsMvKKIgqpfxkj5sNuuSBpiI/w0EGHv7ax5EeTaegpGSqsRWbQdX7X8EpT8IFrTvKeOoMasKAphbaeiA70AxvPJj4dvFYQ440tghHZffFKw0wJEe7yHBuLxtXW+aTgAKobZrea634fSzDivh03Vun8IntqKr+g5I0WT28a6JI23nih+w4c1a2NCW+61GroOYKJEzQx7hNRDKCaK+Ubji/FQlYiMt+4xpdBjRILgnzxAZ/lQ0WH/+3LsG0uQ/DuQlePqChoetkO1IdmusY7VpkMMpjoiqR6G41sMLF3/tO4WSHKxvnw9gGseSaN6QAXon5u+YDWG//7xcayFd1xPZFWkPqaLq2oQVUIrOqfG2n6x9AyS9WRaG18sZWNWjSqvU28lTCOAyrEImI0aX54793OaXczijX5F4oXr2bLDenvmJ8IZE6Vbrlb8JwO7qZ5TVC7HY1U+wcVk9e/6bSnbDYJgaT25uee0x8Z4goyeTvu/aM+MULl1JHEmhI2OgnPCFJPlj1UFPT5ARlR3I4AWdRYH6e7LJE7BK+pWGwOpW6ocYsyZNPGhO1iQY6/VC1pFTApyoA1w4m4JX3UnyZS9Iw51XXT8nhNm/4xL9xL3QPBwmOjl1JxCTMN1/IKBGgxKBkUs8ok5Ktxh1vXtpP+9aPNhEOVSVH7xJiyMzrOZBOuo=-----END CERTIFICATE-----

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
