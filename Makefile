## Variables
REGISTRY := registry.fallow.app
IMAGE := mooz
PLATFORM := linux/arm64
DOCKER_TAG := $(shell cat version.txt)
PLATFORM_TAG := arm

.PHONY: build
phone:
	@echo "Building Docker Image: $(IMAGE):$(DOCKER_TAG)-$(PLATFORM_TAG)..."
	docker buildx build --platform $(PLATFORM) -t $(REGISTRY)/$(IMAGE):$(DOCKER_TAG)-$(PLATFORM_TAG) .

.PHONY: push
push:	
	@echo "Pushing Docker Image: $(IMAGE):$(DOCKER_TAG)-$(PLATFORM_TAG)..."
	docker push $(REGISTRY)/$(IMAGE):$(DOCKER_TAG)-$(PLATFORM_TAG)

.PHONY: docker-run
docker-run:
	docker run -it --rm $(REGISTRY)/$(IMAGE):$(DOCKER_TAG)-$(PLATFORM_TAG)

.PHONY: dev-server
dev-server:
	go run ./cmd/web

.PHONY: release
release: build push
	