FROM golang:1.16-alpine

WORKDIR /app

RUN apk add git


ENV GO111MODULE=on

RUN go mod init github.com/martina97/sdccPROVA
RUN go mod tidy

COPY register/*.go ./register/
COPY register2/*.go ./register2/

RUN go install ./register
RUN go install ./register2