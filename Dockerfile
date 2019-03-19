FROM golang:alpine

RUN mkdir /app
RUN apk update && apk upgrade && \
    apk add --no-cache bash git

ADD . /app
WORKDIR /app/bot
RUN go build -o main .

CMD ["/app/bot/main"]
