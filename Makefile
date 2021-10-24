SHELL := /bin/bash

VERSION := $(shell git rev-parse --short HEAD)
# VERSION := latest

docker-build:
	docker build \
		-f zarf/docker/Dockerfile \
		-t ghcr.io/taraktikos/go-service-amd64:$(VERSION) \
		--build-arg VCS_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

docker-push:
	docker push ghcr.io/taraktikos/go-service:$(VERSION)

compose-up:
	docker compose -f zarf/docker-compose.dev.yml up --build

compose-down:
	docker compose -f zarf/docker-compose.dev.yml down


KIND_CLUSTER := go-service-demo-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yml
	kubectl config set-context --current --namespace=demo-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	cd zarf/k8s/kind/go-service-pod; kustomize edit set image go-service-image=ghcr.io/taraktikos/go-service-amd64:$(VERSION)
	kind load docker-image ghcr.io/taraktikos/go-service-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build zarf/k8s/kind/go-service-pod | kubectl apply -f -

kind-logs:
	kubectl logs -l app=go-service --all-containers=true -f --tail=100

kind-services-delete:
	kustomize build zarf/k8s/kind/go-service-pod | kubectl delete -f -
	kustomize build zarf/k8s/kind/database-pod | kubectl delete -f -

