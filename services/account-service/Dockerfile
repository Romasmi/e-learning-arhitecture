FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy the shared api module
COPY api/ ./api/

# Copy the service go.mod/sum
COPY services/account-service/go.mod services/account-service/go.sum ./services/account-service/

WORKDIR /app/services/account-service
RUN go mod download

# Copy the rest of the service source
COPY services/account-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o account ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/services/account-service/account .
COPY --from=builder --chown=appuser:appgroup /app/services/account-service/migrations ./migrations

USER appuser

EXPOSE 8080 8000
CMD ["./account"]
