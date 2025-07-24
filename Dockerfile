# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache \
	ca-certificates \
	curl

WORKDIR /app

RUN curl -fLs -o /app/deltachat-rpc-server "https://github.com/chatmail/core/releases/download/v2.6.0/deltachat-rpc-server-x86_64-linux" \
	&& chmod +x /app/deltachat-rpc-server


# Copy go mod and sum files
COPY go.mod ./

# Download dependencies first for better caching
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o deltachat-bot .

# Final stage
FROM alpine:latest

WORKDIR /app

ENV PATH=$PATH:/app

# Copy the binary from builder
COPY --from=builder /app/deltachat-rpc-server .
COPY --from=builder /app/deltachat-bot .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./deltachat-bot"]