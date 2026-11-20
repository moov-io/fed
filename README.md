[![Moov Banner Logo](https://user-images.githubusercontent.com/20115216/104214617-885b3c80-53ec-11eb-8ce0-9fc745fb5bfc.png)](https://github.com/moov-io)

<p align="center">
  <a href="https://moov-io.github.io/fed/">Project Documentation</a>
  ·
  <a href="https://moov-io.github.io/fed/api/#overview">API Endpoints</a>
  ·
  <a href="https://moov.io/blog/craft/fed-api-guide/">API Guide</a>
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
[![Apache 2 License](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/moov-io/fed/master/LICENSE)
[![Slack Channel](https://img.shields.io/badge/slack-moov--io-e01563.svg?logo=slack)](https://slack.moov.io/)
[![Docker Pulls](https://img.shields.io/docker/pulls/moov/fed)](https://hub.docker.com/r/moov/fed)
[![GitHub Stars](https://img.shields.io/github/stars/moov-io/fed)](https://github.com/moov-io/fed)
[![Twitter](https://img.shields.io/twitter/follow/moov?style=social)](https://twitter.com/moov?lang=en)

# moov-io/fed

Moov's mission is to give developers an easy way to create and integrate bank processing into their own software products. Our open source projects are each focused on solving a single responsibility in financial services and designed around performance, scalability, and ease of use.

Fed implements utility services for searching the United States Federal Reserve System such as [ABA routing numbers](https://en.wikipedia.org/wiki/ABA_routing_transit_number), financial institution name lookup, and [Fedwire](https://en.wikipedia.org/wiki/Fedwire) and [FedACH](https://en.wikipedia.org/wiki/FedACH) routing information. The HTTP server is available as a [release binary](#binary-distribution), a [Docker image](#docker), or [built from source](#build-from-source), and the Go package `github.com/moov-io/fed` is available. Moov's primary usage for this project is with ACH origination in our [paygate](https://github.com/moov-io/paygate) project.

The data and formats in this repository represent a compilation of **FedWire** and **FedACH** data from the [Federal Reserve Bank Services site](https://frbservices.org/). Both the official Fed plaintext and JSON file formats are supported.

## Table of contents

- [Project status](#project-status)
- [Usage](#usage)
  - [Getting current directory data](#getting-current-directory-data)
  - As an API
    - [Binary](#binary-distribution)
    - [Build from source](#build-from-source)
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

#### Getting current directory data

The files the HTTP server loads are the plaintext (or official FRB JSON) directories in [`data/FedACHdir.txt`](./data/FedACHdir.txt) and [`data/fpddir.txt`](./data/fpddir.txt). The format of those files is documented in [`docs/FedACHdir.md`](./docs/FedACHdir.md) and [`docs/fpddir.md`](./docs/fpddir.md) (and the layout docs they link to). The copies bundled in this repository are from 2018 (latest record dates 2018-11-30 for FedACH and 2018-12-03 for Fedwire). The Federal Reserve Banks no longer publish the complete E-Payments Routing Directory on the public website (as of December 9, 2018) and licensing prevents us from distributing more recent files. The Fed still compiles this data. The server logs a warning at startup when the latest record date is more than a year old.

You can retrieve current files in the official Federal Reserve plaintext or JSON formats and point Fed at them with `FEDACH_DATA_PATH`, `FEDWIRE_DATA_PATH`, or `INITIAL_DATA_DIRECTORY` (see [Configuration](#configuration-settings)).

##### Federal Reserve eServices download code

The Federal Reserve Banks provide the complete E-Payments Routing Directory (FedACH and Fedwire participant files) to credentialed users through FedLine Web and FedLine Advantage. Financial institutions that participate in FedACH or Fedwire can retrieve the files from FedLine or request an E-Payments Routing Directory download code from FedLine to authorize access (including for their corporate customers). Non-banking institutions should contact the financial institution they use for ACH and wire payments; that institution may need to request a download code from FedLine. See the [E-Payments Routing Directory](https://www.frbservices.org/resources/e-payments-routing-directory/) and the [download code request letter template](https://www.frbservices.org/wp-content/uploads/epayments-routing-directory-cover-letter.pdf) published by Federal Reserve Financial Services.

Once you have a routing number (org id) and download code, run Fed with:

```
FRB_ROUTING_NUMBER=123456780
FRB_DOWNLOAD_CODE=86cfa5a9-1ab9-4af5-bd89-0f84d546de13
```

Fed downloads JSON from `https://frbservices.org/EPaymentsDirectory/directories/%s?format=json` (where `%s` is `fedach` or `fedwire`). A successful download has HTTP 200 and an application-level `"response":{"code":100}` with a non-empty participant list. A missing or incorrect download code typically still returns HTTP 200, but the body looks like:

```json
{
  "fedACHParticipants": {
    "response": { "code": 202 },
    "fedACHParticipants": []
  }
}
```

The server logs the HTTP status and FRB response code and exits with an error instead of treating that as a successful download.

Optionally override the download URL (Fed replaces `%s` with `fedach` or `fedwire`):

```
FRB_DOWNLOAD_URL_TEMPLATE=https://my.example.com/files/%s?format=json
```

##### LexisNexis and other vendors

[LexisNexis](https://risk.lexisnexis.com/financial-services/payments-efficiency/payment-routing) and some financial institutions also distribute routing directories. Fed does not implement a vendor-specific parser. It accepts the official Federal Reserve plaintext layout (see the [format docs](./docs/FedACHdir.md)) and the official FRB JSON wrapper (`fedACHParticipants` / `fedwireParticipants` with `response.code`). If a vendor file is in one of those formats, set `FEDACH_DATA_PATH` / `FEDWIRE_DATA_PATH` (or drop the file into `INITIAL_DATA_DIRECTORY` under one of the expected names).

#### Binary distribution

Download the [latest Moov Fed server release](https://github.com/moov-io/fed/releases/latest) (`fed-linux-amd64`, `fed-darwin-amd64`, or `fed.exe`). The server needs FedACH and Fedwire data files; a release binary started from an empty directory exits with `problem opening ./data/FedACHdir.txt`.

Run the binary from a source checkout (so the bundled `data/` files are present), or set data paths / an FRB download code:

```
# from a git checkout
./fed-linux-amd64

# or point at files you have downloaded
FEDACH_DATA_PATH=/path/to/FedACHdir.txt FEDWIRE_DATA_PATH=/path/to/fpddir.txt ./fed-linux-amd64

# or a directory containing one of the expected filenames
INITIAL_DATA_DIRECTORY=/path/to/dir ./fed-linux-amd64

# or download current files from FRB eServices
FRB_ROUTING_NUMBER=... FRB_DOWNLOAD_CODE=... ./fed-linux-amd64
```

The HTTP API listens on `:8086` (admin on `:9096`). See [docs/usage-binary.md](./docs/usage-binary.md) for more examples.

#### Build from source

Requires [Go v1.26 or newer](https://golang.org/doc/install).

```
git clone https://github.com/moov-io/fed.git
cd fed
make build
./bin/server
```

`make build` writes `./bin/server`.

#### Docker

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
      "viewCode": "1",
      "cleanName": "VERIDIAN CREDIT UNION"
    }
  ],
  "stats": {
    "records": 18198,
    "latest": "2018-11-30T00:00:00Z"
  }
}
```

A name search (`?name=Veridian&limit=1`) returns routing number `073983141` in the bundled 2018 data — Veridian has more than one routing number. See [docs/usage-binary.md](./docs/usage-binary.md).

#### **Wire routing number example**

Fed can be used to look up Financial Institutions for [Fedwire](https://en.wikipedia.org/wiki/Fedwire) messages by their routing number (`?routingNumber=...`):

```
curl "localhost:8086/fed/wire/search?routingNumber=273976369"
```
```
{
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
      "date": "20141107",
      "cleanName": "VERIDIAN CREDIT UNION"
    }
  ],
  "stats": {
    "records": 7693,
    "latest": "2018-12-03T00:00:00Z"
  }
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
| `FRB_DOWNLOAD_CODE`         | Federal Reserve Board eServices download code used to download FedACH and FedWire files. A bad code typically returns HTTP 200 with `"response":{"code":202}` and an empty list. | Empty                                                                                                                     |
| `FRB_DOWNLOAD_URL_TEMPLATE` | URL Template for downloading files from alternate source                                              | `https://frbservices.org/EPaymentsDirectory/directories/%s?format=json`                                                   |
| `LOG_FORMAT`                | Format for logging lines to be written as.                                                            | Options: `json`, `plain` - Default: `plain`                                                                               |
| `HTTP_BIND_ADDRESS`         | Address for Fed to bind its HTTP server on. This overrides the command-line flag `-http.addr`.        | Default: `:8086`                                                                                                          |
| `HTTP_ADMIN_BIND_ADDRESS`   | Address for Fed to bind its admin HTTP server on. This overrides the command-line flag `-admin.addr`. | Default: `:9096`                                                                                                          |
| `HTTP_ADMIN_PPROF_ENABLED`  | Set `true`/`1`/`yes` to enable or `false`/`0`/`no` to disable Go pprof on the admin server. Unset keeps historical pprof on. | Unset (enabled)                                                                                                           |
| `HTTP_ADMIN_PPROF_SECRET`   | Optional bearer token required for pprof when it is enabled.                                          | Empty                                                                                                                     |
| `HTTPS_CERT_FILE`           | Filepath containing a certificate intermediate chain to be served by the HTTP server.                 | Empty                                                                                                                     |
| `HTTPS_KEY_FILE`            | Filepath of a private key matching the leaf certificate from `HTTPS_CERT_FILE`.                       | Empty                                                                                                                     |
| `BASE_PATH`                 | HTTP path to serve the API and web UI from.                                                           | `/`                                                                                                                       |

### Data persistence
By design, Fed  **does not persist** (save) any data about the search queries created. The only storage occurs in memory of the process and upon restart Fed will have no files or data saved. Also, no in-memory encryption of the data is performed.

### Go library

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go v1.26 or newer. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/fed/releases/latest) as well. We highly recommend you use a tagged release for production.

```
git clone https://github.com/moov-io/fed.git

# Pull down into the Go Module cache
go get -u github.com/moov-io/fed

go doc github.com/moov-io/fed ACHDictionary
```

Load a directory file and look up a routing number:

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/moov-io/fed"
)

func main() {
	fd, err := os.Open("data/FedACHdir.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	dict := fed.NewACHDictionary()
	if err := dict.Read(fd); err != nil {
		log.Fatal(err)
	}

	p := dict.RoutingNumberSearchSingle("273976369")
	if p == nil {
		log.Fatal("routing number not found")
	}
	fmt.Printf("%s %s\n", p.RoutingNumber, p.CustomerName)
	// 273976369 VERIDIAN CREDIT UNION
}
```

To call a running HTTP server from Go, use the generated client (`github.com/moov-io/fed/client`). See [docs/usage-go.md](./docs/usage-go.md) and [client/README.md](./client/README.md).

## Learn about Fed services participation
- [Intro to Fedwire](https://www.frbservices.org/wp-content/uploads/funds.pdf)
- [Intro to FedACH](https://www.frbservices.org/wp-content/uploads/ach-product-sheet.pdf)
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
- Docker images for `linux/amd64` and `linux/arm64`

Note: 32-bit platforms have known issues and are not supported.

## Contributing

Yes please! Please review our [Contributing guide](https://github.com/moov-io/ach/blob/master/CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) to get started! This repository does not include a `CONTRIBUTING.md`; the ACH guide describes the Moov contribution process.

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go v1.26 or newer. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/fed/releases/latest) as well. We highly recommend you use a tagged release for production.

### Releasing

To make a release of fed simply open a pull request with `CHANGELOG.md` and `version.go` updated with the next version number and details. You'll also need to push the tag (i.e. `git push origin v1.0.0`) to origin in order for CI to make the release.

### Testing

We maintain a comprehensive suite of unit tests and recommend table-driven testing when a particular function warrants several very similar test cases. To run all test files in the current directory, use `go test`. Current releases are built from the `Version` const in `version.go`; `make dist` does not inject the version with ldflags, so bump the const when cutting a release.

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
