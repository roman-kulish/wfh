FROM golang:alpine

ADD . /go/src/github.com/roman-kulish/wfh

RUN go install github.com/roman-kulish/wfh

CMD ["/go/bin/wfh"]

EXPOSE 80