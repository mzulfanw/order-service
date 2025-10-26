# Order Service

Order Service is a microservice for managing orders. This service is built with Go and integrates with PostgreSQL, Redis, and RabbitMQ.

## Features
- Create and retrieve order data
- Integration with Product Service
- Inter-service communication using RabbitMQ
- Caching with Redis

## Folder Structure
- `cmd/` : Application entry point
- `configs/` : Application configuration
- `database/migrations/` : Database migration files
- `internal/` : Business logic (domain, handler, repository, service)
- `pkg/` : External utilities (database, redis, rabbitmq, product client)
- `response/` : Response helpers

## How to Run

### Run the Service
```bash
make run
```

### Database Migration
Make sure you have installed [golang-migrate](https://github.com/golang-migrate/migrate).

```bash
make migrate
```

### Environment Configuration
```
cp .env .env.example
```

## Testing
Run unit tests with:
```bash
go test ./...
```

## API Endpoints
- `POST /orders` : Create a new order
- `GET /orders/product/{productId}` : Get orders by product ID
