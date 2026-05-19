# E-Learning Platform

A microservices architecture demo for a multi-tenant e-learning SaaS. 
Each company gets a branded portal (course catalogue) and an LMS (back-office). 
The project is deployed locally on **minikube** using Helm.

---

## How it works

```
Client → Traefik (API Gateway) → microservice → PostgreSQL / Kafka
                    ↕
              auth-service  (JWT validation on every protected route)
```

**Traefik** is the single entry point. For every protected request it calls `auth-service` 
to validate the bearer token before forwarding to the target service ([authentication flow](docs/sequence_diagram/authentication.puml)).

Each service owns its own PostgreSQL database and publishes domain events to **Redpanda** (Kafka-compatible). 
**Prometheus** scrapes metrics from all services; 
**Grafana** displays per-service dashboards ([metrics flow](docs/sequence_diagram/metrics.puml)).

---

## Services

| Service | Responsibility | Routes |
|---|---|---|
| `auth-service` | Login, JWT issue & validation | `POST /auth/login`, `GET /auth` |
| `portal-service` | Portal & LMS config management | `/portals` |
| `account-service` | Accounts & admin users | `/accounts` |
| `course-service` | Course authoring (chapters, lessons) | `/courses`, `/chapters` |
| `media-service` | Asset upload, transcode jobs | `/media` |
| `student-service` | Students & groups | `/students` |

---

## Stack

| Layer | Technology |
|---|---|
| Services | Go |
| API Gateway | Traefik |
| Database | PostgreSQL (1 primary + 2 read replicas) |
| Message bus | Redpanda (Kafka-compatible) |
| Metrics | Prometheus + Grafana |
| Deployment | Kubernetes / Helm / minikube |
| API contracts | protobuf / gRPC + HTTP gateway |

---

## Quick start

```bash
# 1. Start minikube (8 GB RAM, 4 CPUs recommended)
make minikube-start

# 2. Build images and deploy everything
make up

# 3. In a separate terminal — start port-forward
make run

# API is now available at http://arch.homework:8080
# Traefik dashboard: http://arch.homework:8080/dashboard/  (admin / password)
```

Seed initial data and run tests:

```bash
make test-postman   # Postman collection via Newman
make test-load      # k6 load test
```

Other useful commands:

```bash
make status         # check pod/service status
make grafana-run    # Grafana on http://localhost:3000
make prometheus-run # Prometheus on http://localhost:9090
make restart        # rolling restart of all services
make clean          # remove all Helm releases and namespaces
```

---

## Docs

- [System design](docs/elearning-system-design.md) — full requirements and design decisions
- [C4 diagrams](docs/readme.md) — context, container, component views
- [Deployment diagram](docs/deployment/deployment.puml) — Kubernetes deployment overview
- [Authentication flow](docs/sequence_diagram/authentication.puml) — login and ForwardAuth sequence
- [Metrics flow](docs/sequence_diagram/metrics.puml) — Prometheus pull model
- [Course access check](docs/sequence_diagram/check_access.md) — assignment + license validation
- [Progress tracking](docs/sequence_diagram/progress.md) — student progress sequence
