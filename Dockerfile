FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o server ./cmd/

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/server ./server

CMD ["./server"]
