---
layout: page
title: Binary distribution
hide_hero: true
show_sidebar: false
menubar: docs-menu
---

# Binary distribution

Download the [latest Moov Fed server release](https://github.com/moov-io/fed/releases/latest) for your operating system (`fed-linux-amd64`, `fed-darwin-amd64`, or `fed.exe`).

The server **requires FedACH and Fedwire directory files**. A release binary started from an empty directory exits immediately:

```
problem opening ./data/FedACHdir.txt: open ./data/FedACHdir.txt: no such file or directory
```

Run it from a source checkout (the repository includes `data/FedACHdir.txt` and `data/fpddir.txt`), or set one of the following:

- `FEDACH_DATA_PATH` and `FEDWIRE_DATA_PATH` — paths to the directory files (defaults: `./data/FedACHdir.txt` and `./data/fpddir.txt`)
- `INITIAL_DATA_DIRECTORY` — a directory that contains one of the expected names (`FedACHdir.txt`, `fedachdir.json`, `fedach.txt`, `fedach.json` for ACH; `fpddir.txt`, `fpddir.json`, `fedwire.txt`, `fedwire.json` for Wire)
- `FRB_ROUTING_NUMBER` and `FRB_DOWNLOAD_CODE` — download current files from Federal Reserve eServices (see [Getting current directory data](/#getting-current-directory-data) in the project README)

The bundled files are from 2018. The server logs a warning at startup when the latest record date is more than a year old.

To build from source instead of a release binary, run `make build` in a checkout (requires [Go v1.26 or newer](https://golang.org/doc/install)). That writes `./bin/server`.

```sh
git clone https://github.com/moov-io/fed.git
cd fed
make build
./bin/server
```

```
ts=2026-11-15T16:05:24Z msg="Starting fed server version v0.16.0" level=info
ts=2026-11-15T16:05:24Z msg="search: loading ./data/FedACHdir.txt for ACH data" level=info
ts=2026-11-15T16:05:24Z msg="search: loading ./data/fpddir.txt for Wire data" level=info
ts=2026-11-15T16:05:24Z msg="Finished refresh of ACH data" level=info records=18198
ts=2026-11-15T16:05:24Z msg="Finished refresh of WIRE data" level=info records=7693
ts=2026-11-15T16:05:24Z msg="FedACH data is stale: latest record date is 2018-11-30; replace data/FedACHdir.txt with a current directory file" level=warn
ts=2026-11-15T16:05:24Z msg="Fedwire data is stale: latest record date is 2018-12-03; replace data/fpddir.txt with a current directory file" level=warn
ts=2026-11-15T16:05:24Z msg="binding to :8086 for HTTP server" level=info
```

The Moov Fed service will be running on port `8086` (with an admin port on `9096`). `GET /version` on the admin port reports `v0.16.0`.

Confirm that the service is running by issuing the following command or simply visiting [localhost:8086/ping](http://localhost:8086/ping) in your browser.

```sh
$ curl http://localhost:8086/ping
PONG
```

Search for a routing number:

```
$ curl "localhost:8086/fed/ach/search?routingNumber=273976369"
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

Search for a financial institution by name:

```
$ curl "localhost:8086/fed/ach/search?name=Veridian&limit=1"
{
  "achParticipants": [
    {
      "routingNumber": "073983141",
      "officeCode": "O",
      "servicingFRBNumber": "071000301",
      "recordTypeCode": "1",
      "revised": "122713",
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

Veridian Credit Union has more than one routing number in the bundled 2018 data. A name search with `limit=1` returns `073983141`; the routing-number search above is for `273976369`.

## Fedwire search

```
$ curl "localhost:8086/fed/wire/search?routingNumber=273976369"
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

## FRB eServices download errors

When `FRB_ROUTING_NUMBER` and `FRB_DOWNLOAD_CODE` are set, Fed downloads JSON from Federal Reserve eServices. A successful download has HTTP 200 and an application-level `"response":{"code":100}` with a non-empty participant list.

A missing or incorrect download code typically still returns HTTP 200, but the body looks like:

```json
{
  "fedACHParticipants": {
    "response": { "code": 202 },
    "fedACHParticipants": []
  }
}
```

The server logs the HTTP status and FRB response code and exits with an error (instead of "downloaded ACH file" followed by "read zero records"). See [Getting current directory data](https://github.com/moov-io/fed#getting-current-directory-data) in the project README for how to request a download code.
