.POSIX:
.SUFFIXES:

TEST ?= $$(go list ./... | grep -v 'vendor')
HOSTNAME = hashicorp.com
NAMESPACE = xabinapal
NAME = airflow
BINARY = terraform-provider-${NAME}
VERSION = 0.1.0
OS_ARCH = darwin_amd64

AIRFLOW_VERSION = 2.0.0b2

define APIGEN
	oapi-codegen \
		-generate $(OPTION) \
		-o ./api/$(OPTION).go \
		-package api \
		https://raw.githubusercontent.com/apache/airflow/$(AIRFLOW_VERSION)/airflow/api_connexion/openapi/v1.yaml
endef

default: install

build:
	go build -ldflags $(LD_FLAGS) -o ${BINARY}

release:
	GOOS=darwin  GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386   go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm   go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux   GOARCH=386   go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux   GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux   GOARCH=arm   go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386   go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386   go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

.PHONY: apigen_types
apigen_types: OPTION=types
apigen_types:
	$(APIGEN)

.PHONY: apigen_client
apigen_client: OPTION=client
apigen_client:
	$(APIGEN)

.PHONY: apigen
apigen: apigen_types apigen_client