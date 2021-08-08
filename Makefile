SHELL := /bin/bash

VERSION := $(shell git rev-parse --short HEAD)

build-server:
	docker build \
		-f zarf/docker/Dockerfile \
		-t ghcr.io/taraktikos/go-service:$(VERSION) \
		--build-arg VCS_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.
