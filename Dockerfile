# Build Application
FROM golang:1.22.5-alpine as build
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz

# Run Application
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /app/migrate /app/migrate
COPY app.env .
COPY scripts/start.sh .
COPY db/migration ./db/migration
EXPOSE 8080
ENTRYPOINT [ "/app/start.sh" ]