# Stage 1: Build the Go binary
FROM golang:1.22.3 AS builder

WORKDIR /app

# Copy the Go source code
COPY go.mod go.sum *.go ./

# Download Go modules
RUN go mod download
RUN go test -v

# need to use the CGO_ENABLED=0 to disably dynamic linking used by the net library
RUN CGO_ENABLED=0 GOOS=linux go build -o ./jokes .

# # Stage 2: Create a minimal image with the Go binary
FROM scratch

# # Copy the Go binary from the builder stage
COPY --from=builder /app/jokes /jokes
# # copy public folder to the minimal image
COPY public ./public

# # Expose the application port
EXPOSE 8080

# Command to run the Go binary
CMD ["/jokes"]
