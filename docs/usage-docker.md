---
layout: page
title: Docker
hide_hero: true
show_sidebar: false
menubar: docs-menu
---

# Docker

We publish a [public Docker image `moov/fed`](https://hub.docker.com/r/moov/fed/) from Docker Hub or use this repository. No configuration is required to serve on `:8086` and metrics at `:9096/metrics` in Prometheus format. We also have Docker images for [OpenShift](https://quay.io/repository/moov/fed?tab=tags) published as `quay.io/moov/fed`.

Moov ImageCashLetter is dependent on Docker being properly installed and running on your machine. Ensure that Docker is running. If your Docker client has issues connecting to the service, review the [Docker getting started guide](https://docs.docker.com/get-started/).

```
docker ps
```
```
CONTAINER ID        IMAGE        COMMAND        CREATED        STATUS        PORTS        NAMES
```

Pull & start the Docker image:
```
docker pull moov/fed:latest
docker run -p 8086:8086 -p 9096:9096 moov/fed:latest
```

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