FROM golang:1.12-stretch as builder
WORKDIR /go/src/github.com/moov-io/fed
RUN apt-get update && apt-get install make gcc g++
COPY . .
ENV GO111MODULE=on
RUN go mod download
RUN make build

FROM debian:9
RUN apt-get update && apt-get install -y ca-certificates

# Bring in example data files
COPY data/FedACHdir.txt /data/FedACHdir.txt
COPY data/fpddir.txt /data/fpddir.txt
ENV FEDACH_DATA_PATH=/data/FedACHdir.txt
ENV FEDWIRE_DATA_PATH=/data/fpddir.txt

# Copy server executable
COPY --from=builder /go/src/github.com/moov-io/fed/bin/server /bin/server
# USER moov # TODO(adam): non-root users

EXPOSE 8080
EXPOSE 9090
ENTRYPOINT ["/bin/server"]
