FROM golang:1.13-buster as builder
WORKDIR /go/src/github.com/moov-io/fed
RUN apt-get update && apt-get install make gcc g++
COPY . .
ENV GO111MODULE=on
RUN go mod download
RUN make build
RUN useradd --shell /bin/false moov

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/moov-io/fed/bin/server /bin/server

COPY data/FedACHdir.txt /data/fed/FedACHdir.txt
COPY data/fpddir.txt /data/fed/fpddir.txt

ENV FEDACH_DATA_PATH=/data/fed/FedACHdir.txt
ENV FEDWIRE_DATA_PATH=/data/fed/fpddir.txt

COPY --from=builder /etc/passwd /etc/passwd
USER moov

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/bin/server"]
