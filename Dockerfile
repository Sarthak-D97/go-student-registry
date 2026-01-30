# Stage 1: Build
FROM golang:alpine AS builder

RUN apk add --no-cache gcc musl-dev git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Updated path to ./main.go
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/stuAPI/main.go

# Stage 2: Runtime
FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite

WORKDIR /app

COPY --from=builder /app/main .
COPY config/local.yaml ./config/local.yaml

EXPOSE 8082

CMD ["./main"]