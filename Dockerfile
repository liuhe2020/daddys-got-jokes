From golang:1.22.3 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o daddys-got-jokes

FROM scratch
COPY --from=builder /app/daddys-got-jokes /daddys-got-jokes
EXPOSE 8000
CMD ["/daddys-got-jokes"]