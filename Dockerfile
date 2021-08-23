FROM golang:1.17-buster as builder
WORKDIR /go/src/github.com/moov-io/fed
RUN apt-get update && apt-get install make gcc g++
COPY . .
RUN make build
RUN useradd --shell /bin/false moov

FROM golang:1.17-buster
LABEL maintainer="Moov <support@moov.io>"
RUN apt-get update && apt-get install ca-certificates

COPY --from=builder /go/src/github.com/moov-io/fed/bin/server /bin/server

COPY data/*.txt /data/fed/

ENV FEDACH_DATA_PATH=/data/fed/FedACHdir.txt
ENV FEDWIRE_DATA_PATH=/data/fed/fpddir.txt

COPY --from=builder /etc/passwd /etc/passwd
USER moov

EXPOSE 8086
EXPOSE 9096
ENTRYPOINT ["/bin/server"]
