FROM golang:1.8.1

RUN go get -u \
      github.com/golang/dep/...

WORKDIR /go/src/github.com/minodisk/presigner
COPY . .
RUN go build -o /usr/local/bin/presigner

CMD presigner -help
