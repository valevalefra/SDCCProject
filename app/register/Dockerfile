FROM golang:1.16-alpine

WORKDIR /app

RUN apk add git


ENV GO111MODULE=on


RUN go mod init github.com/valevalefra/SDCCProject
RUN go mod tidy

COPY register/*.go ./register/
COPY utils/*.go /usr/local/go/src/SDCCProject/app/utils/


RUN go install ./register





