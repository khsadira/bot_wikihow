FROM golang:latest

LABEL maintainer="Khan sadirac <khan.sadirac42@gmail.com"

WORKDIR /app

COPY . /app

RUN go build -o wikihow_bot

EXPOSE 9001

CMD ["./wikihow_bot"]
