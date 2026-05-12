DOCKER_IMG := account-service:latest
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
K8S_NAMESPACE := e-learning-system

docker-build:
	docker build \
		--platform linux/amd64 \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f Dockerfile ../..
	@if minikube status >/dev/null 2>&1; then \
		echo "Loading image into minikube..."; \
		minikube image load $(DOCKER_IMG); \
	fi

docker-push:
	docker push $(DOCKER_IMG)

kube-restart:
	kubectl rollout restart deployment/account-service -n $(K8S_NAMESPACE)

kube-status:
	kubectl rollout status deployment/account-service -n $(K8S_NAMESPACE)

kube-logs:
	kubectl logs deployment/account-service -n $(K8S_NAMESPACE)
