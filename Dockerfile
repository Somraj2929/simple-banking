# Stage 1: Build the Go binary
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Create a small image with the binary
FROM alpine:latest
WORKDIR /app/
COPY --from=builder /app/main .
COPY app.env .
RUN chmod +x /app/main
EXPOSE 8080
CMD ["./main"]
