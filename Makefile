HOSTNAME=local
BINARY=terraform-provider-rpt
OS_ARCH=darwin_arm64
NAMESPACE=cmarkulin
NAME=rpt
VERSION=0.54

default: install

build:
	go build -o ${BINARY}

install: build
	mkdir -p ./example/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ./example/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}