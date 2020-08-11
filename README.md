moov-io/fed
===
[![GoDoc](https://godoc.org/github.com/moov-io/fed?status.svg)](https://godoc.org/github.com/moov-io/fed)
[![Build Status](https://github.com/moov-io/fed/workflows/Go/badge.svg)](https://github.com/moov-io/fed/actions)
[![Coverage Status](https://codecov.io/gh/moov-io/fed/branch/master/graph/badge.svg)](https://codecov.io/gh/moov-io/fed)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/fed)](https://goreportcard.com/report/github.com/moov-io/fed)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/fed/master/LICENSE)

Package `github.com/moov-io/fed` implements utility services for searching the United States Federal Reserve System such as [ABA routing numbers](https://en.wikipedia.org/wiki/ABA_routing_transit_number), Financial Institution name lookup and [Fed Wire](https://en.wikipedia.org/wiki/Fedwire) routing information. Moov's primary usage for this project is with ACH origination in our [paygate](https://github.com/moov-io/paygate) project.

Docs: [Project](https://moov-io.github.io/fed/) | [API Endpoints](https://moov-io.github.io/fed/api/)

### Project Status

Moov Fed is actively used in multiple production environments. Please star the project if you are interested in its progress. We would appreciate any issues created or pull requests. Thanks!

### Usage

The `github.com/moov-io/fed` Go package offers search for FEDACH and FEDWIRE Participants.

To get started using Fed download [the latest release](https://github.com/moov-io/fed/releases/latest) or our [Docker image](https://hub.docker.com/r/moov/fed/tags). We also have docker images for [OpenShift](https://quay.io/repository/moov/fed?tab=tags).

**Note**: The Docker image ships with **old data files** (`FedACHdir.txt` and `fpddir.txt`) as example data. In a production deployment updated files should be **obtained from your Financial Institution** and provided to the server process. The official JSON file format from the Federal Reserve is also supported.

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

### Client Library

Fed ships a client library generated from an [OpenAPI Specification](https://en.wikipedia.org/wiki/OpenAPI_Specification) (Go package [`github.com/moov-io/fed/client`](https://godoc.org/github.com/moov-io/fed/client)). We generate a Go version as Moov's primary language is Go, but other languages can be generated. To make a change edit `openapi.yaml` and run `make generate`. Commit the changes and open a pull request.

### Configuration

| Environmental Variable | Description | Default |
|-----|-----|-----|
| `FEDACH_DATA_PATH` | Filepath to FEDACH data file | `./data/FedACHdir.txt` |
| `FEDWIRE_DATA_PATH` | Filepath to FedWIRE data file | `./data/fpddir.txt` |
| `LOG_FORMAT` | Format for logging lines to be written as. | Options: `json`, `plain` - Default: `plain` |
| `HTTP_BIND_ADDRESS` | Address for paygate to bind its HTTP server on. This overrides the command-line flag `-http.addr`. | Default: `:8086` |
| `HTTP_ADMIN_BIND_ADDRESS` | Address for paygate to bind its admin HTTP server on. This overrides the command-line flag `-admin.addr`. | Default: `:9096` |
| `HTTPS_CERT_FILE` | Filepath containing a certificate (or intermediate chain) to be served by the HTTP server. Requires all traffic be over secure HTTP. | Empty |
| `HTTPS_KEY_FILE`  | Filepath of a private key matching the leaf certificate from `HTTPS_CERT_FILE`. | Empty |


## FedWire and FedACH data from the Federal Reserve Bank Services

The data and formats in this repository represent a compilation of the **FedWire** and **FedACH** data from the [Federal Reserve Bank Services site](https://frbservices.org/). Both the official Fed plaintext and JSON file formats are supported.

### FedWire Directory

* [FedWire](./docs/fpddir.md)

### FedACH Directory

* [FedACH](./docs/FedACHdir.md)

### Other resources

* [State and Territory Abbreviations](./docs/Fed_STATE_CODES.md)

## Getting Help

 channel | info
 ------- | -------
[Project Documentation](https://moov-io.github.io/fed/) | Our project documentation available online.
Google Group [moov-users](https://groups.google.com/forum/#!forum/moov-users)| The Moov users Google group is for contributors other people contributing to the Moov project. You can join them without a google account by sending an email to [moov-users+subscribe@googlegroups.com](mailto:moov-users+subscribe@googlegroups.com). After receiving the join-request message, you can simply reply to that to confirm the subscription.
Twitter [@moov_io](https://twitter.com/moov_io)	| You can follow Moov.IO's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channel to have an interactive discussion about the development of the project.

### Copyright and Terms of Use

(c) Federal Reserve Banks

By accessing the [data](./data/) in this repository you agree to the [Federal Reserve Banks' Terms of Use](https://frbservices.org/terms/index.html) and the [E-Payments Routing Directory Terms of Use Agreement](https://www.frbservices.org/EPaymentsDirectory/agreement.html).

## Disclaimer

**THIS REPOSITORY IS NOT AFFILIATED WITH THE FEDERAL RESERVE BANKS AND IS NOT AN OFFICIAL SOURCE FOR THE FEDWIRE AND THE FEDACH DATA.**
