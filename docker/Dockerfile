FROM golang:1.23.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux go build -o /app/web ./cmd/web
RUN GOOS=linux go build -o /app/migrate ./cmd/migrate

FROM golang:1.23.0

WORKDIR /app

COPY --from=builder /app/web /app/web
COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/docker/entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]

