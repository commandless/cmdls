
FROM golang:1.13-alpine as builder
RUN apk add npm
WORKDIR /
COPY . /
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o /go/bin/cmdls
