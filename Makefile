
NAME := push2s3
ARCH := amd64
VERSION := 1.0
DATE := $(shell date)
COMMIT_ID := $(shell git rev-parse --short HEAD)
SDK_INFO := $(shell go version)
LD_FLAGS := -X main.buildInfo 'Version: $(VERSION), commitID: $(COMMIT_ID), build date: $(DATE), SDK: $(SDK_INFO)'

all: clean binaries 

test:
	echo godep go test

binaries: test 
	CGO_ENABLED=0 godep go build -ldflags "$(LD_FLAGS)" -o $(NAME)-darwin-$(ARCH)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 godep go build -ldflags "$(LD_FLAGS)" -o $(NAME)-linux-$(ARCH)

clean: 
	go clean
