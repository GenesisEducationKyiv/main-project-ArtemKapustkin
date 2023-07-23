FROM golang:latest AS build

LABEL maintainer="logs-consumer"

WORKDIR /logs-consumer

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /logs-consumer/main ./

CMD ["./main"]
