
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download



COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .


FROM alpine:latest

WORKDIR /app

# Create a non-root user
RUN apk add --no-cache ca-certificates && addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/main .

# Change ownership of the application binary
RUN chown appuser:appgroup /app/main

# Switch to non-root user
USER appuser

EXPOSE 8080

CMD ["./main"]
