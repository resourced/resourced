FROM golang

ADD . /go/src/github.com/resourced/resourced

RUN mkdir /resourced

ENV RESOURCED_CONFIG_DIR=/go/src/github.com/resourced/resourced/tests/resourced-configs

WORKDIR /go/src/github.com/resourced/resourced

RUN go get ./...
RUN go build
RUN ./resourced

EXPOSE 55555