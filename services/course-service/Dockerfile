FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy the shared api module
COPY api/ ./api/

# Copy the service go.mod/sum
COPY services/course-service/go.mod services/course-service/go.sum ./services/course-service/

WORKDIR /app/services/course-service
RUN go mod download

# Copy the rest of the service source
COPY services/course-service/ .

ARG LDFLAGS="-s -w"
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="${LDFLAGS}" -o course ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/services/course-service/course .
COPY --from=builder --chown=appuser:appgroup /app/services/course-service/migrations ./migrations

USER appuser

EXPOSE 8080 8000
CMD ["./course"]
