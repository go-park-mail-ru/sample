FROM golang:alpine

COPY . /src
WORKDIR /src

RUN apk add --update git gcc musl-dev && GO111MODULE=on go build

EXPOSE 8080

CMD ["./sample"]
