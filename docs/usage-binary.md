---
layout: page
title: Binary distribution
hide_hero: true
show_sidebar: false
menubar: docs-menu
---

# Binary distribution

Download the [latest Moov Fed server release](https://github.com/moov-io/fed/releases) for your operating system and run it from a terminal.

```sh
$ ./fed-darwin-amd64
ts=2019-06-20T23:23:44.870717Z caller=main.go:75 startup="Starting fed server version v0.4.1"
ts=2019-06-20T23:23:44.871623Z caller=main.go:135 transport=HTTP addr=:8086
ts=2019-06-20T23:23:44.871692Z caller=main.go:125 admin="listening on :9096"
```

The Moov Fed service will be running on port `8086` (with an admin port on `9096`).

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
      "viewCode": "1"
    }
  ],
  "wireParticipants": null
}
```

Search for a financial institution by name:

```
$ curl "localhost:8086/fed/ach/search?name=Veridian&limit=1"
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