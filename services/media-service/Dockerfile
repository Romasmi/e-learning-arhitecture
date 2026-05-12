FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy the shared api module
COPY api/ ./api/

# Copy the service go.mod/sum
COPY services/media-service/go.mod services/media-service/go.sum ./services/media-service/

WORKDIR /app/services/media-service
RUN go mod download

# Copy the rest of the service source
COPY services/media-service/ .

ARG LDFLAGS="-s -w"
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="${LDFLAGS}" -o media ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/services/media-service/media .
COPY --from=builder --chown=appuser:appgroup /app/services/media-service/migrations ./migrations

USER appuser

EXPOSE 8080 8000
CMD ["./media"]
