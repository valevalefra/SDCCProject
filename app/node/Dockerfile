FROM golang:1.16-alpine

WORKDIR /app

RUN apk add git


ENV GO111MODULE=on


RUN go mod init github.com/valevalefra/SDCCProject
RUN go mod tidy

COPY node/*.go ./node/
COPY node/sharedFile /node/sharedFile
COPY utility/*.go /usr/local/go/src/SDCCProject/app/utility/



RUN go install ./node

