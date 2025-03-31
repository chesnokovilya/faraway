# Build stage
FROM --platform=linux/arm64 mirror.gcr.io/golang AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the application
RUN  go build -o server ./main.go

# Final stage
FROM --platform=linux/arm64 mirror.gcr.io/golang

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .
# Copy resources
COPY lorem-ipsum.txt .

EXPOSE 8080

CMD ["./server"]