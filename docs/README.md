# Moov Fed

**[GitHub](https://github.com/moov-io/fed)** | **[Running](https://github.com/moov-io/fed#usage)** | **[Configuration](https://github.com/moov-io/fed#configuration-settings)**

## Purpose

[Moov Fed](https://github.com/moov-io/fed) implements an HTTP interface to search [Fedwire](https://github.com/moov-io/fed/tree/master/docs/fpddir.md) and [FedACH](https://github.com/moov-io/fed/tree/master/docs/FedACHdir.md) data from the Federal Reserve Bank Services.

The data and formats below represent a compilation of  **Fedwire** and **FedACH** data from the [Federal Reserve Bank Services site](https://frbservices.org/):

* [FEDACH](https://github.com/moov-io/fed/tree/master/docs/FedACHdir.md)

* [FEDWire](https://github.com/moov-io/fed/tree/master/docs/fpddir.md)

Fed can be used standalone to search for routing numbers by Financial Institution name, city, state, postal code, and routing number. It can also be used in conjunction with [ACH](https://github.com/moov-io/ach) and [WIRE](https://github.com/moov-io/wire) to validate routing numbers.

## Running Moov FED Server

Moov FED can be deployed in multiple scenarios.

- <a href="#binary-distribution">Binary Distributions</a> are released with every versioned release. These are frequently added to the VM/AMI build script for applications needing Moov Fed.
- A <a href="#docker-container">Docker container</a> is built and added to Docker Hub with every versioned release.
- Our hosted [api.moov.io](https://api.moov.io) is updated with every versioned release. Our Kubernetes example is what Moov utilizes in our production environment.

### Binary Distribution

Download the [latest Moov Fed server release](https://github.com/moov-io/fed/releases) for your operating system and run it from a terminal.

```sh
$ cd ~/Downloads/
host:Downloads $ ./fed-darwin-amd64
ts=2019-06-20T23:23:44.870717Z caller=main.go:75 startup="Starting fed server version v0.4.1"
ts=2019-06-20T23:23:44.871623Z caller=main.go:135 transport=HTTP addr=:8086
ts=2019-06-20T23:23:44.871692Z caller=main.go:125 admin="listening on :9096"
```

Next: [Connect to Moov Fed](#connecting-to-moov-fed)

### Docker Container

Moov Fed is dependent on Docker being properly installed and running on your machine. Ensure that Docker is running. If your Docker client has issues connecting to the service, review the [Docker getting started guide](https://docs.docker.com/get-started/).

```sh
$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
```

Execute the Docker run command

```sh
$ docker run -p "8086:8086" -p "9096:9096" moov/fed:latest
ts=2019-06-21T17:03:23.782592Z caller=main.go:69 startup="Starting fed server version v0.4.1"
ts=2019-06-21T17:03:23.78314Z caller=main.go:129 transport=HTTP addr=:8086
ts=2019-06-21T17:03:23.783252Z caller=main.go:119 admin="listening on :9096"
```

Next: [Connect to Moov Fed](#connecting-to-moov-fed)

### Kubernetes

The following snippet runs the Fed Server on [Kubernetes](https://kubernetes.io/docs/tutorials/kubernetes-basics/) in the `apps` namespace. You can reach the fed instance at the following URL from inside the cluster.

```
# Needs to be ran from inside the cluster
$ curl http://fed.apps.svc.cluster.local:8086/ping
PONG
```

Kubernetes manifest - save in a file (`fed.yaml`) and apply with `kubectl apply -f fed.yaml`.

Next: [Connect to Moov Fed](#connecting-to-moov-fed)

## FED Data Files

The data files included in this repository ([`FedACHdir.md`](FedACHdir.md) and [`fpddir.md`](fpddir.md)) are **outdated** and from 2018. The Fed no longer releases this data publicly and licensing on more recent files prevents us from distributing them. However, the Fed still complies this data and you can retrieve up-to-date files for use in our project, either from [Accuity](https://accuity.com/product/us-payments-routing-and-transit-number-file/) or your financial institution.

Moov Fed can read the data files from anywhere on the filesystem. This allows you to mount the files and set `FEDACH_DATA_PATH` / `FEDWIRE_DATA_PATH` environmental variables. Both official formats from the Federal Reserve (plaintext and JSON) are supported.

## Connecting to Moov Fed
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

Search for a Financial Institution by name:

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

### API Documentation

See our [API documentation](https://moov-io.github.io/fed/api/) for Moov Fed endpoints.


### Other Resources

- [State and Territory Abbreviations](./Fed_STATE_CODES.md)
- [U.S. Department of the Treasury FAQ](https://www.treasury.gov/resource-center/faqs/Sanctions/Pages/faq_general.aspx#basic)

### Copyright and Terms of Use

Copyright &copy; Federal Reserve Banks

By accessing the [data](https://github.com/moov-io/fed/tree/master/data) in this repository you agree to the [Federal Reserve Banks' Terms of Use](https://frbservices.org/terms/index.html) and the [E-Payments Routing Directory Terms of Use Agreement](https://www.frbservices.org/EPaymentsDirectory/agreement.html).

## Disclaimer

**THIS REPOSITORY IS NOT AFFILIATED WITH THE FEDERAL RESERVE BANKS AND IS NOT AN OFFICIAL SOURCE FOR FEDWIRE AND FEDACH DATA.**

## Getting Help

 channel | info
 ------- | -------
 [Project Documentation](https://moov-io.github.io/fed/) | Our project documentation available online.
Twitter [@moov_io](https://twitter.com/moov_io)	| You can follow Moov.IO's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io/fed/issues) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channel (`#fed`) to have an interactive discussion about the development of the project.
