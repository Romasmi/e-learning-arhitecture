FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copy the shared api module
COPY api/ ./api/

# Copy the service go.mod/sum
COPY services/student-service/go.mod services/student-service/go.sum ./services/student-service/

WORKDIR /app/services/student-service
RUN go mod download

# Copy the rest of the service source
COPY services/student-service/ .

ARG LDFLAGS="-s -w"
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="${LDFLAGS}" -o student ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app
RUN chown -R appuser:appgroup /app

COPY --from=builder --chown=appuser:appgroup /app/services/student-service/student .
COPY --from=builder --chown=appuser:appgroup /app/services/student-service/migrations ./migrations

USER appuser

EXPOSE 9000 8000
CMD ["./student"]
