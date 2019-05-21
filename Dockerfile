FROM golang:alpine as builder

RUN mkdir /build
ADD . /build
WORKDIR /build

RUN apk update && apk upgrade && \
    apk add --no-cache bash git

RUN go build -o bin/main ./bot

FROM alpine
COPY --from=builder /build/bin /app/bin
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/app/bin/main"]
