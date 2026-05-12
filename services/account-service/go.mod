module github.com/elearning/account-service

go 1.25.1

replace github.com/Romasmi/e-learning-arhitecture => ../../api

require (
	github.com/Romasmi/e-learning-arhitecture v0.0.0-00010101000000-000000000000
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0
	github.com/jackc/pgx/v5 v5.9.2
	github.com/joho/godotenv v1.5.1
	github.com/prometheus/client_golang v1.20.5
	github.com/segmentio/kafka-go v0.4.51
	go.uber.org/zap v1.28.0
	google.golang.org/grpc v1.81.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.16 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260504160031-60b97b32f348 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260427160629-7cedc36a6bc4 // indirect
)
