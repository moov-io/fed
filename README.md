moov-io/fed
===
[![GoDoc](https://godoc.org/github.com/moov-io/fed?status.svg)](https://godoc.org/github.com/moov-io/fed)
[![Build Status](https://travis-ci.com/moov-io/fed.svg?branch=master)](https://travis-ci.com/moov-io/fed)
[![Coverage Status](https://codecov.io/gh/moov-io/fed/branch/master/graph/badge.svg)](https://codecov.io/gh/moov-io/fed)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/fed)](https://goreportcard.com/report/github.com/moov-io/fed)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/fed/master/LICENSE)

*project is under active development and is not production ready*

Package `github.com/moov-io/fed` implements utility services for searching the United States Federal Reserve System.

### Configuration

**Search Similarity Metrics**

FED computes string similarity using the JaroWinkler and Levenshtein algorithm and can match sensitivity can be configured with environment variables.

| Environmental Variable | Description | Default |
|-----|-----|-----|
| `ACHJaroWinklerSimilarity` | Ratio of boosting the score of exact matches at the beginning of the strings. | 0.85 |
| `ACHLevenshteinSimilarity` | Ratio of Levenshtein distance for two strings to be considered equal. | 0.85 |

### Usage

Go library
github.com/moov-io/fed offers a Go based search for FEDACH and FEDWIRE Participants.

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

By accesing the [data](./data/) in this repository you agree to the [Federal Reserve Banks' Terms of Use](https://frbservices.org/terms/index.html) and the [E-Payments Routing Directory Terms of Use Agreement](https://www.frbservices.org/EPaymentsDirectory/agreement.html).  

## Disclaimer

**THIS REPOSITORY IS NOT AFFILIATED WITH THE FEDERAL RESERVE BANKS AND IS NOT AN OFFICIAL SOURCE FOR THE FEDWIRE AND THE FEDACH DATA.**
