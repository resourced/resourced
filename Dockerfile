FROM golang

ADD . /go/src/github.com/resourced/resourced

RUN mkdir /resourced

ENV RESOURCED_DB /resourced/db
ENV RESOURCED_CONFIG_READER_DIR=/go/src/github.com/resourced/resourced/tests/data/config-reader
ENV RESOURCED_CONFIG_WRITER_DIR=/go/src/github.com/resourced/resourced/tests/data/config-writer

WORKDIR /go/src/github.com/resourced/resourced

RUN go get ./...
RUN go build
RUN ./resourced

EXPOSE 55555