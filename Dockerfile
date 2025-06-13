FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOGACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o server ./cmd/

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/server ./server

CMD ["./server"]
