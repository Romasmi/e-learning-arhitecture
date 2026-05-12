# API Skill

This project uses an **API-First** approach with Protocol Buffers (gRPC) and gRPC Gateway (REST).

## 1. Proto Definitions
- Source files: `api/*.proto`.
- External dependencies: `google/api/*.proto`.
- *Rule*: Always use `option go_package = "github.com/Romasmi/e-learning-arhitecture/billing-service/api;api";`.

## 2. Code Generation
- Tool: `buf`.
- Command: `make generate`.
- Output: `api/` (for Go code) and `api/swagger/` (for Swagger JSON).
- *Rule*: Never edit files in `api/` manually.

## 3. Serving Docs
- **Swagger UI**: Available at `/swagger/`. It renders the JSON files from `api/swagger/`.
- **Protos**: Source `.proto` files are served at `/proto/`.

## 4. REST Mapping
- Use `google.api.http` options in `.proto` files to define REST endpoints.
- The gRPC Gateway automatically generates the HTTP reverse proxy.
