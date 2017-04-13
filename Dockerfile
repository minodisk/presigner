FROM golang:1.8.1

RUN go get -u \
      github.com/golang/dep/...

WORKDIR /go/src/github.com/minodisk/presigner
COPY . .
RUN go build -o /usr/local/bin/presigner

CMD presigner \
      -id $GOOGLE_AUTH_EMAIL \
      -key $GOOGLE_AUTH_KEY \
      -bucket $PRESIGNER_BUCKET \
      -duration $PRESIGNER_DURATION \
      -port $PRESIGNER_PORT
