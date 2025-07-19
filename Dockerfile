FROM golang:1.24-alpine as builder

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g cmd/server/main.go

RUN go build -o auth-service cmd/server/main.go

FROM alpine:latest

WORKDIR /app

RUN mkdir -p /app/configs
COPY configs/config.yaml /app/configs/auth-service.yaml

COPY --from=builder /app/auth-service /app/

EXPOSE 8080

CMD ["/app/auth-service"]