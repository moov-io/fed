---
layout: page
title: Overview
hide_hero: true
show_sidebar: false
menubar: docs-menu
---

# Overview

![Moov Fed Logo](https://user-images.githubusercontent.com/20115216/107704449-176ca600-6c72-11eb-8963-45bd9ae3e557.jpg)

Moov's mission is to give developers an easy way to create and integrate bank processing into their own software products. Our open source projects are each focused on solving a single responsibility in financial services and designed around performance, scalability, and ease of use.

Fed implements utility services for searching the United States Federal Reserve System such as [ABA routing numbers](https://en.wikipedia.org/wiki/ABA_routing_transit_number), financial institution name lookup, and [Fedwire](https://en.wikipedia.org/wiki/Fedwire) and [FedACH](https://en.wikipedia.org/wiki/FedACH) routing information. The HTTP server is available in a [Docker image](#docker) and the Go package `github.com/moov-io/fed` is available. Moov's primary usage for this project is with ACH origination in our [paygate](https://github.com/moov-io/paygate) project.

The data and formats in this repository represent a compilation of **FedWire** and **FedACH** data from the [Federal Reserve Bank Services site](https://frbservices.org/). Both the official Fed plaintext and JSON file formats are supported.

## Copyright and terms of use

Copyright &copy; Federal Reserve Banks

By accessing the [data](https://github.com/moov-io/fed/tree/master/data) in this repository you agree to the [Federal Reserve Banks' Terms of Use](https://frbservices.org/terms/index.html) and the [E-Payments Routing Directory Terms of Use Agreement](https://www.frbservices.org/EPaymentsDirectory/agreement.html).

## Disclaimer

**THIS REPOSITORY IS NOT AFFILIATED WITH THE FEDERAL RESERVE BANKS AND IS NOT AN OFFICIAL SOURCE FOR FEDWIRE AND FEDACH DATA.**