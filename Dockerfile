# Build Application
FROM golang:1.22.5-alpine as build
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Application
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE 8080
CMD ["/app/main"]