FROM golang:1.12-alpine

WORKDIR /app/src/notification
COPY . /app/src/notification

ENV GOPATH /app

RUN apk add git \
    && go get -u github.com/golang/dep/cmd/dep \
    && /app/bin/dep ensure \
    && go build -o build/http internal/cmd/http.go \
    && rm -rf /var/cache/apk/*

CMD ["build/http"]