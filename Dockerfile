FROM alpine:latest

COPY build/linux/elaho /

CMD ["/elaho"]