moov-io/fed
===
[![GoDoc](https://godoc.org/github.com/moov-io/fed?status.svg)](https://godoc.org/github.com/moov-io/fed)
[![Build Status](https://travis-ci.com/moov-io/fed.svg?branch=master)](https://travis-ci.com/moov-io/fed)
[![Coverage Status](https://codecov.io/gh/moov-io/fed/branch/master/graph/badge.svg)](https://codecov.io/gh/moov-io/fed)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/fed)](https://goreportcard.com/report/github.com/moov-io/fed)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/fed/master/LICENSE)

Package `github.com/moov-io/fed` implements utility services for searching the United States Federal Reserve System such as [ABA routing numbers](https://en.wikipedia.org/wiki/ABA_routing_transit_number), Financial Institution name lookup and [Fed Wire](https://en.wikipedia.org/wiki/Fedwire) routing information. Moov's primary usage for this project is with ACH origination in our [paygate](https://github.com/moov-io/paygate) project.

Docs: [docs.moov.io](https://docs.moov.io/en/latest/) | [api docs](https://api.moov.io/apps/fed/)

### Project Status

Moov FED is under active development and in production for multiple companies. Please star the project if you are interested in its progress.

### Usage

Go library
github.com/moov-io/fed offers a Go based search for FEDACH and FEDWIRE Participants.

To get started using Fed download [the latest release](https://github.com/moov-io/fed/releases) or our [Docker image](https://hub.docker.com/r/moov/fed/tags).

Docs: [docs.moov.io](https://docs.moov.io/en/latest/) | [api docs](https://api.moov.io/apps/fed/)

Note: The Docker image ships with old data files (`FedACHdir.txt` and `fpddir.txt`) as example data. In a production deployment updated files should be obtained from your Financial Institution and provided to the server process.

#### ACH Routing Number Example

Fed can be used to lookup a Financial Institutions for Automated Clearing House ([ACH](https://en.wikipedia.org/wiki/Automated_Clearing_House)) transfers by their routing number (`?routingNumber=...`) or name (`?name=...`):

```
$ curl -s localhost:8086/fed/ach/search?routingNumber=273976369 | jq .
{
  "achParticipants": [
    {
      "routingNumber": "273976369",
      "officeCode": "O",
      "servicingFRBNumber": "071000301",
      "recordTypeCode": "1",
      "revised": "041513",
      "newRoutingNumber": "000000000",
      "customerName": "VERIDIAN CREDIT UNION",
      "achLocation": {
        "address": "1827 ANSBOROUGH",
        "city": "WATERLOO",
        "state": "IA",
        "postalCode": "50702",
        "postalCodeExtension": "0000"
      },
      "phoneNumber": "3192878332",
      "statusCode": "1",
      "viewCode": "1"
    }
  ],
  "wireParticipants": null
}
```

#### Wire Routing Number Example

Fed can be used to lookup a Financial Institutions for FED Wire Messages ([FEDWire](https://en.wikipedia.org/wiki/Fedwire)) by their routing number (`?routingNumber=...`) or name (`?name=...`):

```
$ curl -s localhost:8086/fed/wire/search?routingNumber=273976369 | jq .
{
  "achParticipants": null,
  "wireParticipants": [
    {
      "routingNumber": "273976369",
      "telegraphicName": "VERIDIAN",
      "customerName": "VERIDIAN CREDIT UNION",
      "wireLocation": {
        "city": "WATERLOO",
        "state": "IA"
      },
      "fundsTransferStatus": "Y",
      "fundsSettlementOnlyStatus": " ",
      "bookEntrySecuritiesTransferStatus": "N",
      "date": "20141107"
    }
  ]
}
```

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
