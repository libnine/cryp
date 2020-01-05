FROM golang:1.13.5 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cryp ./cmd/cryp

FROM alpine:latest AS prod

COPY --from=builder /app .

EXPOSE 8000

CMD ["./cryp"]
