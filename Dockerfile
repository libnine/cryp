FROM golang:1.13.5-alpine3.11
RUN mkdir /cryp
ADD . /cryp
WORKDIR /cryp
RUN go build -o ./cryp/cmd/main ./cryp/bin
CMD ["./cryp/bin/main"]
