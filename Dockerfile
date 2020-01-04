FROM golang:1.13.5-alpine3.11
RUN mkdir /cryp && apk update && apk install git
ADD . /cryp
WORKDIR /cryp
RUN go get -d -v ./... && \
  go install -v ./... && \
  GOOS=linux GOARCH=amd64 go build -o ./bin/main ./cmd/cryp
EXPOSE 8000
CMD ["./bin/main"]
