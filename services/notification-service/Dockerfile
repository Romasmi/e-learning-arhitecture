FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy the shared api module
COPY api/ ./api/

# Copy the service go.mod/sum
COPY services/notification-service/go.mod services/notification-service/go.sum ./services/notification-service/

WORKDIR /app/services/notification-service
RUN go mod download

# Copy the rest of the service source
COPY services/notification-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api    ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o worker ./cmd/worker

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/services/notification-service/api    .
COPY --from=builder --chown=appuser:appgroup /app/services/notification-service/worker .
COPY --from=builder --chown=appuser:appgroup /app/services/notification-service/migrations ./migrations
COPY --from=builder --chown=appuser:appgroup /app/services/notification-service/config.yaml .
USER appuser

EXPOSE 8000
CMD ["./api"]