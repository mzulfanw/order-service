run:
	go run ./cmd/server.go

migrate:
	set -a; source .env; set +a; migrate -path ./database/migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up