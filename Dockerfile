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

EXPOSE 4000

ENTRYPOINT ["/app/web"]

