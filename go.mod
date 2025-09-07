module github.com/ialekseychuk/my-place-identity

go 1.24.6

require (
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/crypto v0.39.0
	google.golang.org/grpc v1.75.0
	google.golang.org/protobuf v1.36.8
)

require (
	github.com/ialekseychuk/my-place-proto v0.0.0-20250906202310-c4bb361e439e // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
)

// replace github.com/ialekseychuk/my-place-proto => ../proto
