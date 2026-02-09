# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

# Production stage
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server .

RUN mkdir -p /app/uploads

VOLUME ["/app/uploads"]

EXPOSE 4000

CMD ["./server"]
