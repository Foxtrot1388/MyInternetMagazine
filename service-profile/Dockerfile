FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /usr/local/go/src/
ADD go.mod .
ADD go.sum .
RUN go mod download

RUN apk --no-cache add git gcc libc-dev

ADD . .
RUN CGO_ENABLED=1 go build -tags musl -mod=mod -o app.exe cmd/app/main.go

#lightweight docker container with binary
FROM alpine:latest

ARG POSTGRES_HOST=$POSTGRES_HOST
ARG POSTGRES_PORT=$POSTGRES_PORT
ARG POSTGRES_USER=$POSTGRES_USER
ARG POSTGRES_PASSWORD=$POSTGRES_PASSWORD
ARG SIGNING_KEY=$SIGNING_KEY
ARG KAFKA_HOST=$KAFKA_HOST

RUN apk --no-cache add ca-certificates

COPY --from=builder /usr/local/go/src/app.exe /
COPY /migrations /migrations

EXPOSE 8081
EXPOSE 8082

CMD [ "/app.exe"]