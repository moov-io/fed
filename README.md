[![Moov Banner Logo](https://user-images.githubusercontent.com/20115216/104214617-885b3c80-53ec-11eb-8ce0-9fc745fb5bfc.png)](https://github.com/moov-io)

<p align="center">
  <a href="https://moov-io.github.io/fed/">Project Documentation</a>
  ·
  <a href="https://moov-io.github.io/fed/api/#overview">API Endpoints</a>
  ·
  <a href="https://moov.io/blog/education/fed-api-guide/">API Guide</a>
  ·
  <a href="https://slack.moov.io/">Community</a>
  ·
  <a href="https://moov.io/blog/">Blog</a>
  <br>
  <br>
</p>

[![GoDoc](https://godoc.org/github.com/moov-io/fed?status.svg)](https://godoc.org/github.com/moov-io/fed)
[![Build Status](https://github.com/moov-io/fed/workflows/Go/badge.svg)](https://github.com/moov-io/fed/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/fed)](https://goreportcard.com/report/github.com/moov-io/fed)
[![Repo Size](https://img.shields.io/github/languages/code-size/moov-io/fed?label=project%20size)](https://github.com/moov-io/fed)
[![Apache 2 License](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/ach/master/LICENSE)
[![Slack Channel](https://slack.moov.io/badge.svg?bg=e01563&fgColor=fffff)](https://slack.moov.io/)
[![Docker Pulls](https://img.shields.io/docker/pulls/moov/fed)](https://hub.docker.com/r/moov/fed)
[![GitHub Stars](https://img.shields.io/github/stars/moov-io/fed)](https://github.com/moov-io/fed)
[![Twitter](https://img.shields.io/twitter/follow/moov?style=social)](https://twitter.com/moov?lang=en)

# moov-io/fed

Moov's mission is to give developers an easy way to create and integrate bank processing into their own software products. Our open source projects are each focused on solving a single responsibility in financial services and designed around performance, scalability, and ease of use.

Fed implements utility services for searching the United States Federal Reserve System such as [ABA routing numbers](https://en.wikipedia.org/wiki/ABA_routing_transit_number), financial institution name lookup, and [Fedwire](https://en.wikipedia.org/wiki/Fedwire) and [FedACH](https://en.wikipedia.org/wiki/FedACH) routing information. The HTTP server is available in a [Docker image](#docker) and the Go package `github.com/moov-io/fed` is available. Moov's primary usage for this project is with ACH origination in our [paygate](https://github.com/moov-io/paygate) project.

The data and formats in this repository represent a compilation of **FedWire** and **FedACH** data from the [Federal Reserve Bank Services site](https://frbservices.org/). Both the official Fed plaintext and JSON file formats are supported.

## Table of contents

- [Project status](#project-status)
- [Usage](#usage)
  - As an API
    - [Docker](#docker) ([Config](#configuration-settings))
    - [Web UI](#webui)
    - [Google Cloud](#google-cloud-run) ([Config](#configuration-settings))
    - [Data persistence](#data-persistence)
  - [As a Go module](#go-library)
- [Learn about Fed services participation](#learn-about-fed-services-participation)
- [Getting help](#getting-help)
- [Supported and tested platforms](#supported-and-tested-platforms)
- [Contributing](#contributing)
- [Related projects](#related-projects)
- [Copyright](#copyright-and-terms-of-use)

### Project status

Moov Fed is actively used in multiple production environments. Please star the project if you are interested in its progress. We would appreciate any issues created or pull requests. Thanks!

### Usage

The Fed project implements an HTTP server and [Go library](https://pkg.go.dev/github.com/moov-io/fed) for searching for FedACH and Fedwire participants.

**Note**: The data files included in this repository ([`FedACHdir.md`](./docs/FedACHdir.md) and [`fpddir.md`](./docs/fpddir.md)) are **outdated** and from 2018. The Fed no longer releases this data publicly and licensing on more recent files prevents us from distributing them. However, the Fed still complies this data and you can retrieve up-to-date files for use in our project, either from [LexisNexis](https://risk.lexisnexis.com/financial-services/payments-efficiency/payment-routing) or your financial institution.

Moov Fed can read the data files from anywhere on the filesystem. This allows you to mount the files and set `FEDACH_DATA_PATH` / `FEDWIRE_DATA_PATH` environmental variables. Both official formats from the Federal Reserve (plaintext and JSON) are supported.

#### Download files

The Federal Reserve Board (FRB) eServices offers API access to download the files. To download these files, work with your ODFI / banking partner to obtain a download code. Then run Fed with the following environment variables set.

```
FRB_ROUTING_NUMBER=123456780
FRB_DOWNLOAD_CODE=86cfa5a9-1ab9-4af5-bd89-0f84d546de13
```

#### Download files from proxy

Fed can download the files from a proxy or other HTTP resources. The optional URL template is configured as an environment variable. If the URL template is not configured, Fed will download the files directly from FRB eServices by default. This value is considered a template because when preparing the request Fed replaces `%s` in the path with the requested list name(`fedach` or `fedwire`).

```
FRB_DOWNLOAD_URL_TEMPLATE=https://my.example.com/files/%s?format=json
```

### Docker

We publish a [public Docker image `moov/fed`](https://hub.docker.com/r/moov/fed/) from Docker Hub or use this repository. No configuration is required to serve on `:8086` and metrics at `:9096/metrics` in Prometheus format. We also have Docker images for [OpenShift](https://quay.io/repository/moov/fed?tab=tags) published as `quay.io/moov/fed`.

Pull & start the Docker image:
```
docker pull moov/fed:latest
docker run -p 8086:8086 -p 9096:9096 moov/fed:latest
```

### WebUI

With the release of `v0.14.1` Moov Fed contains a builtin Web UI that's hosted from the Go binary / Docker image.

#### **ACH routing number example**

Fed can be used to look up Financial Institutions for Automated Clearing House ([ACH](https://en.wikipedia.org/wiki/Automated_Clearing_House)) transfers by their routing number (`?routingNumber=...`):

```
curl "localhost:8086/fed/ach/search?routingNumber=273976369"
```
```
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

#### **Wire routing number example**

Fed can be used to look up Financial Institutions for [Fedwire](https://en.wikipedia.org/wiki/Fedwire) messages by their routing number (`?routingNumber=...`):

```
curl "localhost:8086/fed/wire/search?routingNumber=273976369"
```
```
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

### Google Cloud Run

To get started in a hosted environment you can deploy this project to the Google Cloud Platform.

From your [Google Cloud dashboard](https://console.cloud.google.com/home/dashboard) create a new project and call it:
```
moov-fed-demo
```

Enable the [Container Registry](https://cloud.google.com/container-registry) API for your project and associate a [billing account](https://cloud.google.com/billing/docs/how-to/manage-billing-account) if needed. Then, open the Cloud Shell terminal and run the following Docker commands, substituting your unique project ID:

```
docker pull moov/fed
docker tag moov/fed gcr.io/<PROJECT-ID>/fed
docker push gcr.io/<PROJECT-ID>/fed
```

Deploy the container to Cloud Run:
```
gcloud run deploy --image gcr.io/<PROJECT-ID>/fed --port 8086
```

Select your target platform to `1`, service name to `fed`, and region to the one closest to you (enable Google API service if a prompt appears). Upon a successful build you will be given a URL where the API has been deployed:

```
https://YOUR-FED-APP-URL.a.run.app
```

Now you can ping the server:
```
curl https://YOUR-FED-APP-URL.a.run.app/ping
```
You should get this response:
```
PONG
```

### Configuration settings

| Environmental Variable      | Description                                                                                           | Default                                                                                                                   |
|-----------------------------|-------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------|
| `FEDACH_DATA_PATH`          | Filepath to FedACH data file                                                                          | `./data/FedACHdir.txt`                                                                                                    |
| `FEDWIRE_DATA_PATH`         | Filepath to Fedwire data file                                                                         | `./data/fpddir.txt`                                                                                                       |
| `INITIAL_DATA_DIRECTORY`    | Directory of files to be used instead of downloading or `*_DATA_PATH` variables.                      | ACH: FedACHdir.txt, fedachdir.json, fedach.txt, fedach.json<br />Wire: fpddir.json, fpddir.txt, fedwire.txt, fedwire.json |
| `FRB_ROUTING_NUMBER`        | Federal Reserve Board eServices (ABA) routing number used to download FedACH and FedWire files        | Empty                                                                                                                     |
| `FRB_DOWNLOAD_CODE`         | Federal Reserve Board eServices (ABA) download code used to download FedACH and FedWire files         | Empty                                                                                                                     |
| `FRB_DOWNLOAD_URL_TEMPLATE` | URL Template for downloading files from alternate source                                              | `https://frbservices.org/EPaymentsDirectory/directories/%s?format=json`                                                   |
| `LOG_FORMAT`                | Format for logging lines to be written as.                                                            | Options: `json`, `plain` - Default: `plain`                                                                               |
| `HTTP_BIND_ADDRESS`         | Address for Fed to bind its HTTP server on. This overrides the command-line flag `-http.addr`.        | Default: `:8086`                                                                                                          |
| `HTTP_ADMIN_BIND_ADDRESS`   | Address for Fed to bind its admin HTTP server on. This overrides the command-line flag `-admin.addr`. | Default: `:9096`                                                                                                          |
| `HTTPS_CERT_FILE`           | Filepath containing a certificate intermediate chain to be served by the HTTP server.                 | Empty                                                                                                                     |
| `HTTPS_KEY_FILE`            | Filepath of a private key matching the leaf certificate from `HTTPS_CERT_FILE`.                       | Empty                                                                                                                     |

### Data persistence
By design, Fed  **does not persist** (save) any data about the search queries created. The only storage occurs in memory of the process and upon restart Fed will have no files or data saved. Also, no in-memory encryption of the data is performed.

### Go library

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go v1.18 or newer. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/fed/releases/latest) as well. We highly recommend you use a tagged release for production.

```
$ git@github.com:moov-io/fed.git

# Pull down into the Go Module cache
$ go get -u github.com/moov-io/fed

$ go doc github.com/moov-io/fed ACHDictionary
```

## Learn about Fed services participation
- [Intro to Fedwire](https://www.frbservices.org/assets/financial-services/wires/funds.pdf)
- [Intro to FedACH](https://www.frbservices.org/assets/financial-services/ach/ach-product-sheet.pdf)
- [U.S. Department of the Treasury FAQ](https://www.treasury.gov/resource-center/faqs/Sanctions/Pages/faq_general.aspx#basic)
- [State and Territory Abbreviations](./docs/Fed_STATE_CODES.md)
- [Fedwire Directory File Format](./docs/fpddir_FORMAT.md)
- [FedACH Directory File Format](./docs/FEDACHDIR_FORMAT.MD)

## Getting help

 channel | info
 ------- | -------
[Project Documentation](https://moov-io.github.io/fed/) | Our project documentation available online.
Twitter [@moov](https://twitter.com/moov)	| You can follow Moov.io's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io/fed/issues) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channel to have an interactive discussion about the development of the project.

## Supported and tested platforms

- 64-bit Linux (Ubuntu, Debian), macOS, and Windows

Note: 32-bit platforms have known issues and are not supported.

## Contributing

Yes please! Please review our [Contributing guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) to get started!

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go v1.18 or newer. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/fed/releases/latest) as well. We highly recommend you use a tagged release for production.

### Releasing

To make a release of fed simply open a pull request with `CHANGELOG.md` and `version.go` updated with the next version number and details. You'll also need to push the tag (i.e. `git push origin v1.0.0`) to origin in order for CI to make the release.

### Testing

We maintain a comprehensive suite of unit tests and recommend table-driven testing when a particular function warrants several very similar test cases. To run all test files in the current directory, use `go test`.

## Related projects
As part of Moov's initiative to offer open source fintech infrastructure, we have a large collection of active projects you may find useful:

- [Moov Watchman](https://github.com/moov-io/watchman) offers search functions over numerous trade sanction lists from the United States and European Union.

- [Moov Image Cash Letter](https://github.com/moov-io/imagecashletter) implements Image Cash Letter (ICL) files used for Check21, X.9 or check truncation files for exchange and remote deposit in the U.S.

- [Moov Wire](https://github.com/moov-io/wire) implements an interface to write files for the Fedwire Funds Service, a real-time gross settlement funds transfer system operated by the United States Federal Reserve Banks.

- [Moov ACH](https://github.com/moov-io/ach) provides ACH file generation and parsing, supporting all Standard Entry Codes for the primary method of money movement throughout the United States.

- [Moov Metro 2](https://github.com/moov-io/metro2) provides a way to easily read, create, and validate Metro 2 format, which is used for consumer credit history reporting by the United States credit bureaus.

## Copyright and terms of use

(c) Federal Reserve Banks

By accessing the [data](./data/) in this repository you agree to the [Federal Reserve Banks' Terms of Use](https://frbservices.org/terms/index.html) and the [E-Payments Routing Directory Terms of Use Agreement](https://www.frbservices.org/EPaymentsDirectory/agreement.html).

## Disclaimer

**THIS REPOSITORY IS NOT AFFILIATED WITH THE FEDERAL RESERVE BANKS AND IS NOT AN OFFICIAL SOURCE FOR FEDWIRE AND FEDACH DATA.**

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
