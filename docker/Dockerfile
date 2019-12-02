FROM alpine

WORKDIR /rayseep

ADD . ./

RUN chmod +x /rayseep/start.sh && mkdir /etc/rayseep

EXPOSE 4301 4302 4303

CMD ["./ray-seep","-c","/etc/rayseep/rayseep.toml"]