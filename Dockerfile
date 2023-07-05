FROM golang:1.18.2-alpine as builder
RUN apk add alpine-sdk
WORKDIR /go/app

RUN apk add git

COPY go.mod /go/app
COPY go.sum /go/app
RUN go mod download

COPY . /go/app


RUN GOOS=linux GOARCH=amd64 go build -o rest-api -tags musl

FROM alpine:latest as runner
WORKDIR /root/
COPY --from=builder /go/app/rest-api .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#RUN ls /root/rest-api
ENTRYPOINT /root/rest-api