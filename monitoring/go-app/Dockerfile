# Use official Golang base image
FROM golang:1.21-alpine as builder

WORKDIR /app
COPY main.go .

# Download necessary dependencies (no modules needed in this case)
RUN go build -o app main.go

# Minimal runtime image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .

# Set binary as entrypoint
ENTRYPOINT ["./app"]
