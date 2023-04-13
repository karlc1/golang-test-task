FROM golang:1.19 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o server ./cmd/reporter/main.go

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server

ENV PORT $API_PORT
EXPOSE $PORT

CMD ["/app/server"]
