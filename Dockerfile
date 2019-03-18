FROM golang:latest AS builder

RUN mkdir /build
ADD . /build
WORKDIR /build
RUN go get -d ./...
RUN go build -o bin/main ./bot

FROM ubuntu:latest
COPY --from=builder /build/bin /app
CMD ["/app/main"]
