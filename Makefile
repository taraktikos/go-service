SHELL := /bin/bash

VERSION := $(shell git rev-parse --short HEAD)

build-server:
	docker build \
		-f zarf/docker/Dockerfile \
		-t bankets-app-amd64:$(VERSION) \
		--build-arg VCS_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.
