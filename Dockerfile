FROM golang:1.13.5 AS builder
LABEL autodelete="true"
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cryp ./cmd/cryp

FROM alpine:latest AS runtime
LABEL description="prod"
COPY --from=builder /app .
EXPOSE 8000
CMD ["./cryp"]