# Stage 1: Build the Go binary
FROM golang:1.22.4 AS builder

WORKDIR /app

# Copy the Go source code
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal

# Download Go modules
RUN go mod download

# need to use the CGO_ENABLED=0 to disably dynamic linking used by the net library
RUN CGO_ENABLED=0 GOOS=linux go build -o ./jokes ./cmd/api

# # Stage 2: Create a minimal image with the Go binary
FROM scratch

# # Copy the Go binary from the builder stage
COPY --from=builder /app/jokes /jokes

# # Expose the application port
EXPOSE 8080

# Command to run the Go binary
CMD ["/jokes"]
