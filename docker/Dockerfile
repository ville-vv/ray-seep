FROM alpine

WORKDIR /go

ADD . ./

RUN chmod +x /go/start.sh

EXPOSE 32202 32201 32203

CMD ["./ray-seep","-c","/etc/rayseep.toml"]