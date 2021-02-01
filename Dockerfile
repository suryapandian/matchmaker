
FROM golang:1.15-alpine AS builder

ADD . /go/src/matchmaker

WORKDIR /go/src/matchmaker

RUN go build -mod=vendor -o matchmaker .

EXPOSE 3001

ENTRYPOINT [ "./matchmaker"]