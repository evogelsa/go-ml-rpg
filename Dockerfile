FROM golang:1.13-buster

MAINTAINER Ethan Vogelsang "evogelsa@iu.edu"

WORKDIR /go/src/go-ml-rpg

COPY src/ ./

RUN go mod download

RUN go build -o /go-ml-rpg

CMD [ "/go-ml-rpg" ]
