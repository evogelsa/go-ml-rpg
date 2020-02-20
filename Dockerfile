FROM golang:1.13-buster

MAINTAINER Ethan Vogelsang "evogelsa@iu.edu"

WORKDIR /go/src/go-ml-rpg

ENTRYPOINT ["/go/src/go-ml-rpg/autorun.sh"]

CMD ["bash"]
