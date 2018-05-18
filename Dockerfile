FROM golang:latest
MAINTAINER Georg V. Sattler <georg@sattler.io>


RUN mkdir -p /go/src/channelManager.sattler.io
WORKDIR /go/src/channelManager.sattler.io

COPY . /go/src/channelManager.sattler.io

EXPOSE 9000

RUN go get
RUN go build -o channelManager .

CMD ["./channelManager"]