FROM golang:1.16-alpine

WORKDIR /app

RUN apk add git


ENV GO111MODULE=on

RUN go mod init github.com/valevalefra/SDCCProject
RUN go mod tidy

COPY register/*.go ./register/
COPY utils/*.go ./utils/
COPY node/*.go ./node/
COPY prova/*.go ./prova/


RUN go install ./register
RUN go install ./node
