GO := go
DOCKER := docker
TAG?=$(shell git rev-parse HEAD)
REGISTRY?=probablynot
IMAGE=evio-scratch
GOOS=linux

all: build

build:
	@echo ">> building using docker"
	@$(DOCKER) build -t ${REGISTRY}/${IMAGE}:${TAG} -f Dockerfile .
	@$(DOCKER) tag ${REGISTRY}/${IMAGE}:${TAG} ${REGISTRY}/${IMAGE}:latest

push:
	docker push ${REGISTRY}/${IMAGE}:${TAG}
	docker push ${REGISTRY}/${IMAGE}:latest

bin:
	@echo ">> building local binary"
	CGO_ENABLED=0 GOOS=${GOOS} go build -a -ldflags '-extldflags "-static"' -o evio-scratch

.PHONY: all build push bin
