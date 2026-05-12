FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy the shared api module
COPY api/ ./api/

# Copy the service go.mod/sum
COPY services/auth-service/go.mod services/auth-service/go.sum ./services/auth-service/

WORKDIR /app/services/auth-service
RUN go mod download

# Copy the rest of the service source
COPY services/auth-service/ .

ARG LDFLAGS="-s -w"
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="${LDFLAGS}" -o auth ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="${LDFLAGS}" -o seed ./cmd/seed

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/services/auth-service/auth .
COPY --from=builder --chown=appuser:appgroup /app/services/auth-service/seed .
COPY --from=builder --chown=appuser:appgroup /app/services/auth-service/migrations ./migrations
COPY --from=builder --chown=appuser:appgroup /app/services/auth-service/config.yaml .
USER appuser

EXPOSE 8000
CMD ["./auth"]
