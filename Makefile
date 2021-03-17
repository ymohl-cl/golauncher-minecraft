BINARY=golauncher-minecraft
APP=launcher
BIN_FOLDER=bin
IGNORED_FOLDER=.ignore
MODULE_NAME := $(shell go list -m)
CONFIG_FILE=config.json

all: install tools lint build

.PHONY: install
install:
	@go mod download

.PHONY: build
build:
	@CGO_ENABLED=1 CC=gcc GOOS=darwin GOARCH=amd64 go build -tags static -ldflags "-s -w" -o ${BIN_FOLDER}/${BINARY} ${MODULE_NAME}/cmd/${APP}
	@cp ${CONFIG_FILE} ${BIN_FOLDER}/${CONFIG_FILE}

test:
	@go test -count=1 ./...

.PHONY: lint
lint:
	golint ./...

.PHONY: tools
tools:
	go get -u golang.org/x/lint/golint

.PHONY: clean
clean:
	@rm -rf ${IGNORED_FOLDER} 

.PHONY: fclean
fclean: clean
	@rm -rf ${BIN_FOLDER}
