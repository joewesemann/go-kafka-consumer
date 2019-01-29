FROM golang:1.8

RUN mkdir -p /go/src/go-kafka
WORKDIR /go/src/go-kafka

ADD . /go/src/go-kafka

RUN go get -v
