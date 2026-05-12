FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy the shared api module
COPY api/ ./api/

# Copy the service go.mod/sum
COPY services/user-service/go.mod services/user-service/go.sum ./services/user-service/

WORKDIR /app/services/user-service
RUN go mod download

# Copy the rest of the service source
COPY services/user-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/services/user-service/main .
COPY --from=builder --chown=appuser:appgroup /app/services/user-service/migrations ./migrations
COPY --from=builder --chown=appuser:appgroup /app/services/user-service/config.yaml .
USER appuser

EXPOSE 8000
CMD ["./main"]