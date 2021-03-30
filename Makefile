BINARY=golauncher-minecraft
APP=launcher
BIN_FOLDER=bin
IGNORED_FOLDER=.ignore
MODULE_NAME := $(shell go list -m)
CONFIG_FILE=config.json

all: install lint test build

.PHONY: install
install:
	@go mod download

.PHONY: build1
build1:
	@CGO_ENABLED=1 CC=gcc GOOS=darwin GOARCH=amd64 go build -buildmode=pie -tags static -ldflags "-s -w" -o ${BIN_FOLDER}/${BINARY} ${MODULE_NAME}/cmd/${APP}

.PHONY: build2
build2:
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -buildmode=pie -tags static -ldflags "-s -w" -o ${BIN_FOLDER}/${BINARY}.exe ${MODULE_NAME}/cmd/${APP}

.PHONY: test
test:
	@go test -count=1 -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: lint
lint:
	@golint -set_exit_status ./...

.PHONY: tools
tools:
	@go get -u golang.org/x/lint/golint

.PHONY: clean
clean:
	@rm -rf ${IGNORED_FOLDER} 

.PHONY: fclean
fclean: clean
	@rm -rf ${BIN_FOLDER}
