FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /usr/local/go/src/
ADD go.mod .
ADD go.sum .
RUN go mod download

RUN apk --no-cache add git gcc libc-dev

ADD . .
RUN CGO_ENABLED=1 go build -tags musl -mod=mod -o app.exe main.go

#lightweight docker container with binary
FROM alpine:latest

ARG KAFKA_HOST=$KAFKA_HOST
ARG KAFKA_TOPIC_EMAIL=$KAFKA_TOPIC_EMAIL
ARG SENDER_EMAIL=$SENDER_EMAIL
ARG PASS_EMAIL=$PASS_EMAIL
ARG SMTP_AUTH_ADDRESS_EMAIL=$SMTP_AUTH_ADDRESS_EMAIL
ARG SMTP_SERVER_EMAIL=$SMTP_SERVER_EMAIL

RUN apk --no-cache add ca-certificates

COPY --from=builder /usr/local/go/src/app.exe /

CMD [ "/app.exe"]