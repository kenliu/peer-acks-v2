FROM golang:1.19
ENV workdir /build
WORKDIR $workdir

RUN go install -v .

#VOLUME ["/data"]
#WORKDIR /data
CMD ["peer-acks-v2"]
