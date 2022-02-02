# Image URL to use all building/pushing image targets
IMG ?= msm-k8s-svc-helper:latest

.PHONY: binary
binary:
	go build -o bin/main ./cmd

# Build the docker image
docker-build:
	docker build . -t ${IMG}

deploy:
	kubectl apply -f deployment/rbac.yaml
	kubectl apply -f deployment/service.yaml
	kubectl apply -f deployment/deployment.yaml

clean-deploy:
	cd deployment && ./destroy.sh

.PHONY: tidy
tidy: ## Execute go mod tidy
	go mod tidy