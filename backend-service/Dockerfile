
# Build stage
FROM golang:alpine AS builder

WORKDIR /build

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o ./out .

# Final stage
FROM alpine:latest

WORKDIR /main

# Copy the binary from builder
COPY --from=builder /build/out ./out

# Make it executable
RUN chmod +x ./out

EXPOSE 3000

# Run the binary
CMD ["./out"]
