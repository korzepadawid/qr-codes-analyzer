# Build stage
FROM golang:1.18.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.14
WORKDIR /app
RUN mkdir -p token/keys
COPY --from=builder /app/token/keys/* ./token/keys
COPY --from=builder /app/main .
EXPOSE 8080
CMD [ "/app/main" ]
