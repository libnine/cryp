FROM golang:1.13.5

RUN mkdir /cryp

ADD . /cryp

WORKDIR /cryp

RUN apk update && \
apk add git && \
go get -d -v ./... && \
go install -v ./... && \
GOOS=linux GOARCH=amd64 go build -o ./bin/main ./cmd/cryp

FROM alpine

EXPOSE 8000

CMD ["./bin/main"]
