FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git curl


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN git clone https://github.com/golang-migrate/migrate.git /tmp/migrate && \
    cd /tmp/migrate/cmd/migrate && \
    go build -tags 'postgres' -o /usr/local/bin/migrate

# Vuild
RUN go build -o app ./cmd/app/main.go

# Stage 2 — Final runtime image
FROM alpine:latest

WORKDIR /root/

# Certs for https requests
RUN apk --no-cache add ca-certificates

# Copy necessary files
COPY --from=builder /app/app .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations ./migrations
COPY envs/.env.dev ./envs/.env.dev

# ENV
ENV DATABASE_URL="postgres://user:password@postgres:5432/app_db?sslmode=disable"
ENV HTTP_PORT=3000

# Run migrations, then start app
CMD /bin/sh -c '\
  echo "Waiting for PostgreSQL..."; \
  for i in $(seq 1 10); do \
    migrate -path=./migrations -database=$DATABASE_URL up && break || echo "Retrying migration..."; \
    sleep 3; \
  done; \
  echo "Starting app..."; \
  ./app'