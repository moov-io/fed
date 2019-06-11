moov-io/fed
===
[![GoDoc](https://godoc.org/github.com/moov-io/fed?status.svg)](https://godoc.org/github.com/moov-io/fed)
[![Build Status](https://travis-ci.com/moov-io/fed.svg?branch=master)](https://travis-ci.com/moov-io/fed)
[![Coverage Status](https://codecov.io/gh/moov-io/fed/branch/master/graph/badge.svg)](https://codecov.io/gh/moov-io/fed)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/fed)](https://goreportcard.com/report/github.com/moov-io/fed)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/fed/master/LICENSE)

*project is under active development and is not production ready*

Package `github.com/moov-io/fed` implements utility services for searching the United States Federal Reserve System such as [ABA routing numbers](https://en.wikipedia.org/wiki/ABA_routing_transit_number), Financial Institution name lookup and [Fed Wire](https://en.wikipedia.org/wiki/Fedwire) routing information. Moov's primary usage for this project is with ACH origination in our [paygate](https://github.com/moov-io/paygate) project.

### Usage

Go library
github.com/moov-io/fed offers a Go based search for FEDACH and FEDWIRE Participants.

To get started using Fed download [the latest release](https://github.com/moov-io/fed/releases) or our [Docker image](https://hub.docker.com/r/moov/fed/tags).

Docs: [docs.moov.io](https://docs.moov.io/en/latest/) | [api docs](https://api.moov.io/apps/fed/)

Note: The Docker image ships with old data files (`FedACHdir.txt` and `fpddir.txt`) as example data. In a production deployment updated files should be obtained from your Financial Institution and provided to the server process.

### Configuration

| Environmental Variable | Description | Default |
|-----|-----|-----|
| `FEDACH_DATA_PATH` | Filepath to FEDACH data file | `./data/FedACHdir.txt` |
| `FEDWIRE_DATA_PATH` | Filepath to FedWIRE data file | `./data/fpddir.txt` |

## FedWire and FedACH data from the Federal Reserve Bank Services

The data and formats in this repository represent a compilation of the **FedWire** and **FedACH** data from the [Federal Reserve Bank Services site](https://frbservices.org/).

### FedWire Directory

* [FedWire](./docs/fpddir.md)

### FedACH Directory

* [FedACH](./docs/FedACHdir.md)

### Other resources

* [State and Territory Abbreviations](./docs/Fed_STATE_CODES.md)

### Copyright and Terms of Use

(c) Federal Reserve Banks

By accessing the [data](./data/) in this repository you agree to the [Federal Reserve Banks' Terms of Use](https://frbservices.org/terms/index.html) and the [E-Payments Routing Directory Terms of Use Agreement](https://www.frbservices.org/EPaymentsDirectory/agreement.html).

## Disclaimer

**THIS REPOSITORY IS NOT AFFILIATED WITH THE FEDERAL RESERVE BANKS AND IS NOT AN OFFICIAL SOURCE FOR THE FEDWIRE AND THE FEDACH DATA.**
