.PHONY: generate build run-api run-worker run-cli compose-up compose-down migrate-up migrate-down



build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go
#	go build -o bin/cli cmd/cli/main.go

run-api:
	go run cmd/api/main.go

run-worker:
	go run cmd/worker/main.go

#run-cli:
#	go run cmd/cli/main.go

compose-up:
	docker-compose -f deployments/docker-compose.yaml up -d

compose-down:
	docker-compose -f deployments/docker-compose.yaml down

migrate-up:
	migrate -path migrations -database "postgres://user:password@localhost:5432/notification_service?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://user:password@localhost:5432/notification_service?sslmode=disable" down

docker-build:
	docker build --platform linux/amd64 -t romasmi/notification-service:latest -f Dockerfile ../..
	@if minikube status >/dev/null 2>&1; then \
		echo "Loading image into minikube..."; \
		minikube image load romasmi/notification-service:latest; \
	fi

docker-push:
	docker push romasmi/notification-service:latest

kube-restart:
	kubectl rollout restart deployment/notification-service -n e-learning-system

kube-status:
	kubectl rollout status deployment/notification-service -n e-learning-system

lint:
	golangci-lint run ./...