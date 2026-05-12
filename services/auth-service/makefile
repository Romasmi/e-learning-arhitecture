DOCKER_IMG := auth-service:latest
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
K8S_NAMESPACE := e-learning-system

docker-build:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f deployments/calendar.Dockerfile .
	@if minikube status >/dev/null 2>&1; then \
		echo "Loading image into minikube..."; \
		minikube image load $(DOCKER_IMG); \
	fi

kube-restart:
	kubectl rollout restart deployment/auth-service -n $(K8S_NAMESPACE)

kube-status:
	kubectl rollout status deployment/auth-service -n $(K8S_NAMESPACE)

kube-logs:
	kubectl logs deployment/auth-service -n $(K8S_NAMESPACE)
