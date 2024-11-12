FROM golang:1.23.3-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -C cmd -o ../go-test-task

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/go-test-task .
COPY app-config.yml .

EXPOSE 8080

CMD ["./go-test-task"]