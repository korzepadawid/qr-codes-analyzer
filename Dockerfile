# Build stage
FROM golang:1.18.5-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-386.tar.gz | tar xvz
RUN go build -o main main.go

# Run stage
FROM alpine:3.14
WORKDIR /app
RUN mkdir -p token/keys
RUN mkdir templates
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migrations ./migration
COPY --from=builder /app/migrate .
COPY --from=builder /app/templates/* ./templates
COPY --from=builder /app/token/keys/* ./token/keys
COPY --from=builder /app/main .
EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
