FROM golang:1.22-alpine AS builder

WORKDIR /app

# Download Go modules
COPY go.mod go.sum* ./
RUN go mod download

# Copy the source code
COPY . .

# Build
RUN CGO_ENABLED=0 go build -o /disky-api ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /disky-api .

EXPOSE 8080

CMD ["./disky-api"]
