FROM golang:1.19
ENV workdir /build
WORKDIR $workdir
COPY . .

RUN go install -v .

#VOLUME ["/data"]
#WORKDIR /data
CMD ["peer-acks-v2"]
