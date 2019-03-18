FROM golang:alpine

RUN mkdir /app
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

ADD . /app
WORKDIR /app/bot
RUN go get -d ./...
RUN go build -o main .

CMD ["/app/bot/main"]
