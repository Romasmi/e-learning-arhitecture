.PHONY: up build deploy restart install-traefik install-db install-kafka install-grafana hosts run wait-db wait-kafka wait-api clean redeploy status help prometheus-run grafana-run forward-kafka forward-db forward-traefik proto-gen

# Main target to start everything from scratch
up: proto-gen build deploy wait-db wait-kafka wait-api

# Proto generation
proto-gen:
	$(MAKE) -C ./api generate

# Build Docker images and load them into minikube
build:
	$(MAKE) -j7 -C ./services/auth-service docker-build & \
	$(MAKE) -j7 -C ./services/portal-service docker-build & \
	$(MAKE) -j7 -C ./services/account-service docker-build & \
	$(MAKE) -j7 -C ./services/media-service docker-build & \
	$(MAKE) -j7 -C ./services/course-service docker-build & \
	$(MAKE) -j7 -C ./services/student-service docker-build & \
	wait

deploy: install-traefik install-db install-kafka install-prometheus install-grafana install-app

install-app:
	helm upgrade --install e-learning-system ./deployments/helm/e-learning-system \
		--namespace e-learning-system \
		--create-namespace

uninstall-app:
	helm uninstall e-learning-system -n e-learning-system --ignore-not-found

# Helm installations
install-traefik:
	helm repo add traefik https://traefik.github.io/charts || true
	helm repo update traefik

	helm upgrade --install e-learning-traefik traefik/traefik \
	  --namespace e-learning-traefik \
	  --create-namespace \
	  -f ./deployments/helm/traefik-values.yaml

forward-traefik:
	kubectl port-forward -n e-learning-traefik $$(kubectl get pods -n e-learning-traefik -o name) 9000:9000

install-db:
	helm repo add bitnami https://repo.broadcom.com/bitnami-files/
	helm repo update bitnami
	helm upgrade --install postgresql bitnami/postgresql \
		--namespace e-learning-system \
		--create-namespace \
		--values deployments/helm/postgresql-values.yaml

db-connect:
	kubectl exec -it postgresql-0 -n e-learning-system -- \
		psql -U user -d postgres

forward-db:
	kubectl port-forward svc/postgresql 5432:5432 -n e-learning-system

install-kafka:
	helm repo add redpanda https://charts.redpanda.com
	helm repo update redpanda
	helm upgrade --install redpanda redpanda/redpanda \
		--namespace e-learning-system \
		--create-namespace \
		--values deployments/helm/redpanda-values.yaml

forward-kafka:
	kubectl port-forward svc/kafka 9092:9092 -n e-learning-system

install-prometheus:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	helm repo update prometheus-community
	helm upgrade --install prometheus prometheus-community/prometheus \
		--namespace e-learning-system \
		--create-namespace

install-grafana:
	helm repo add grafana https://grafana.github.io/helm-charts || true
	helm repo update grafana
	helm upgrade --install grafana grafana/grafana \
		--namespace e-learning-system \
		--create-namespace \
		--values deployments/helm/grafana-values.yaml

# Utility targets
hosts:
	@echo "Updating /etc/hosts for arch.homework..."
	@sudo sed -i '' '/arch.homework/d' /etc/hosts || sudo sed -i '/arch.homework/d' /etc/hosts
	@echo "127.0.0.1 arch.homework" | sudo tee -a /etc/hosts

run:
	@echo "Starting port-forwarding for Traefik... (Keep this running)"
	@echo "Access API: http://arch.homework:8080"
	@echo "Access Traefik Dashboard: http://arch.homework:8080/dashboard/"
	@echo "Access Grafana: http://localhost:3000 (after make grafana-run)"
	kubectl port-forward service/e-learning-traefik 8080:8080 -n e-learning-traefik

wait-db:
	@echo "Waiting for PostgreSQL to be ready..."
	kubectl wait --namespace e-learning-system --for=condition=ready pod -l app.kubernetes.io/name=postgresql --timeout=120s

wait-kafka:
	@echo "Waiting for Redpanda to be ready..."
	kubectl wait --namespace e-learning-system --for=condition=ready pod -l app.kubernetes.io/name=redpanda --timeout=300s

wait-api:
	@echo "Waiting for API deployments to be ready..."
	kubectl rollout status deployment/auth-service -n e-learning-system --timeout=120s
	kubectl rollout status deployment/portal-service -n e-learning-system --timeout=120s
	kubectl rollout status deployment/account-service -n e-learning-system --timeout=120s
	kubectl rollout status deployment/media-service -n e-learning-system --timeout=120s
	kubectl rollout status deployment/course-service -n e-learning-system --timeout=120s
	kubectl rollout status deployment/student-service -n e-learning-system --timeout=120s

restart:
	kubectl rollout restart deployment/auth-service -n e-learning-system
	kubectl rollout restart deployment/portal-service -n e-learning-system
	kubectl rollout restart deployment/account-service -n e-learning-system
	kubectl rollout restart deployment/media-service -n e-learning-system
	kubectl rollout restart deployment/course-service -n e-learning-system
	kubectl rollout restart deployment/student-service -n e-learning-system
	$(MAKE) wait-api

clean:
	helm uninstall e-learning-system -n e-learning-system --ignore-not-found
	helm uninstall e-learning-traefik -n e-learning-traefik --ignore-not-found
	helm uninstall postgresql -n e-learning-system --ignore-not-found
	helm uninstall redpanda -n e-learning-system --ignore-not-found
	helm uninstall prometheus -n e-learning-system --ignore-not-found
	helm uninstall grafana -n e-learning-system --ignore-not-found
	kubectl delete namespace e-learning-traefik --ignore-not-found=true
	kubectl delete namespace e-learning-system --ignore-not-found=true

status:
	@echo "\n--- Infrastructure ---"
	@echo "e-learning-traefik:"
	@kubectl get pods -n e-learning-traefik
	@echo "PostgreSQL:"
	@kubectl get pods -n e-learning-system -l app.kubernetes.io/name=postgresql
	@echo "Redpanda:"
	@kubectl get pods -n e-learning-system -l app.kubernetes.io/name=redpanda
	@echo "\n--- Application ---"
	@kubectl get pods -n e-learning-system -l app=auth-service
	@kubectl get pods -n e-learning-system -l app=portal-service
	@kubectl get pods -n e-learning-system -l app=account-service
	@kubectl get pods -n e-learning-system -l app=media-service
	@kubectl get pods -n e-learning-system -l app=course-service
	@kubectl get pods -n e-learning-system -l app=student-service
	@echo "\n--- Services ---"
	@kubectl get svc -n e-learning-system
	@kubectl get svc -n e-learning-traefik
	@echo "\n--- Routes ---"
	@kubectl get ingressroute -n e-learning-system

prometheus-run:
	kubectl port-forward service/prometheus-server 9090:80 -n e-learning-system

grafana-run:
	kubectl port-forward service/grafana 3000:80 -n e-learning-system

grafana-pass:
	@kubectl get secret grafana -o jsonpath="{.data.admin-password}" -n e-learning-system | base64 --decode ; echo ""

redeploy: build docker-push restart

help:
	@echo "Usage:"
	@echo "  make up          - Build images and deploy everything (from scratch)"
	@echo "  make redeploy    - Build, push and restart all services"
	@echo "  make run         - Start minikube tunnel (required for access)"
	@echo "  make status      - Check deployment status"
	@echo "  make clean       - Remove all resources"
	@echo "  make install-app - Install application using Helm"
	@echo "  make forward-kafka - Port-forward Kafka to localhost:9092"
	@echo "  make forward-db    - Port-forward PostgreSQL to localhost:5432"
	@echo ""
	@echo "Quick Start:"
	@echo "  1. make up"
	@echo "  2. make run (in another terminal)"
	@echo "  3. Access API: http://arch.homework:8080/accounts"
	@echo "  4. Access Dashboard: http://arch.homework:8080/dashboard/"

draw-puml:
	plantuml -tsvg ./docs/puml/*.puml

test-postman: seed-supervisor
	bash tests/postman/run.sh

seed-supervisor:
	$(MAKE) -C services/auth-service kube-seed