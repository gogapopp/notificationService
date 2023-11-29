FROM golang:1.21.1 AS builder

COPY . /github.com/gogapopp/notificationService/
COPY .env /github.com/gogapopp/notificationService/.env
WORKDIR /github.com/gogapopp/notificationService/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/gogapopp/notificationService/.bin/app .
COPY --from=builder /github.com/gogapopp/notificationService/.env ../.env

CMD ["./app"]