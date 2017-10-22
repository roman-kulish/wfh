FROM golang:alpine

ARG BUCKET="https://storage.googleapis.com/wfh/%s.jpg"

ADD . /go/src/github.com/roman-kulish/wfh

RUN go install -ldflags "-X main.bucket=${BUCKET}" github.com/roman-kulish/wfh

CMD ["/go/bin/wfh"]

EXPOSE 8080
