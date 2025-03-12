FROM golang:1.23 AS builder

WORKDIR /app

ENV DATABASE_URL="postgres://postgres:password@postgres:5432/moneydb?sslmode=disable"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest  

WORKDIR /root/

COPY --from=builder /app/main .

# Copy .env file
COPY .env . 

# Expose port
EXPOSE 8080

# Run the app
CMD ["./main"]
