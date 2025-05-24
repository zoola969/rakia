# Stage 1: Build the application
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .


# Stage 2: Create a minimal image
FROM alpine:3.19
WORKDIR /app

COPY --from=builder /main /app/main
COPY blog_data.json /app/

ENTRYPOINT ["./main"]
