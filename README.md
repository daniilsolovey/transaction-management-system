# 🎰 Transaction Management System

This is a high-performance transaction processing service built in Go. It uses Kafka for event streaming, PostgreSQL for persistent storage. It supports ingestion of user transactions like `bet` and `win` through an HTTP API or via Kafka consumers.

---

## ✨ Features

- REST API with [Gin](https://github.com/gin-gonic/gin)
- Kafka producer/consumer integration
- PostgreSQL persistence with upsert logic
- Redis integration for analytics (optional)
- Swagger documentation
- Graceful shutdown, structured logging
- Dockerized setup with Kafka, Zookeeper, PostgreSQL, Redis

---

## 📦 Technologies

| Layer     | Tech                                   |
|-----------|----------------------------------------|
| API       | Gin, Swagger                           |
| Queue     | Kafka + Zookeeper                      |
| DB        | PostgreSQL                             |
| Cache     | Redis                                  |
| Build     | Go 1.21+, Docker                       |
| Logging   | `log/slog`                             |

---

## 🚀 Getting Started

### Prerequisites

- [Docker](https://www.docker.com/)
- [Go 1.21+](https://golang.org/dl/)

---

### 📁 Folder Structure

```
.
├── cmd/                # Main app entry
├── configs/            # Viper-based configuration
├── internal/
│   ├── delivery/
│   │   ├── http/       # Gin handlers
│   │   └── kafka/      # Kafka consumer/producer
│   ├── repository/     # PostgreSQL and Redis implementation
│   ├── usecase/        # Business logic
│   └── domain/         # DTOs and entity definitions
├── docs/               # Swagger auto-generated
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## 🐳 Run with Docker

```bash
docker compose up --build
```

Then access:

- 🧠 Swagger UI: [http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html)
- 📡 Kafka broker: `kafka:9092`
- 🗃️ PostgreSQL: `postgres:5432`
- 🔁 Redis: `redis:6379`

---

## 🔌 API Endpoints

### POST `/transactions`

Create a new transaction (enqueue to Kafka)

**Request Body**:

```json
{
  "user_id": "11111111-1111-1111-1111-111111111111",
  "type": "bet",
  "amount": 120.5
}
```

**Response**:
- `202 Accepted` on success
- `400 Bad Request` if invalid
- `500 Internal Server Error` on Kafka failure

---

### GET `/transactions`

Fetch user transactions.

**Query Params**:
- `user_id` (required)
- `type` = `bet` | `win` (optional)

**Response**:
```json
[
  {
    "user_id": "11111111-1111-1111-1111-111111111111",
    "type": "bet",
    "amount": 100.0,
    "timestamp": "2025-07-21T10:00:00Z"
  }
]
```

---

## 🔄 Kafka Topics

| Topic        | Description           |
|--------------|-----------------------|
| transactions | Incoming transactions |

---

## ⚙️ Environment Variables (`.env.dev`)

```dotenv
APP_ENV=dev
SERVICE_NAME=transaction-manager
HTTP_PORT=3000

DATABASE_URL=postgres://user:password@postgres:5432/app_db?sslmode=disable
DB_MAX_CONNS=5
DB_MAX_CONN_LIFETIME=300s

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

KAFKA_ADDRESS=kafka:9092
KAFKA_TOPIC=transactions
KAFKA_GROUP_ID=transaction-group
```

---

## 🧪 Testing

To test the API manually:

```bash
curl -X POST http://localhost:3000/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "11111111-1111-1111-1111-111111111111",
    "type": "bet",
    "amount": 100.50
  }'
```

Run unit tests:

```bash
make test
```

---

## 🛠️ Development

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
go run cmd/main.go
```


## RUN

```bash
make up
```