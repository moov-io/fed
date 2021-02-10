---
layout: page
title: Overview
hide_hero: true
show_sidebar: false
menubar: docs-menu
---

# Overview

![Moov Fed Logo](https://moov.io/images/social/moov-share-image.png)

Moov's mission is to give developers an easy way to create and integrate bank processing into their own software products. Our open source projects are each focused on solving a single responsibility in financial services and designed around performance, scalability, and ease of use.

Fed implements utility services for searching the United States Federal Reserve System such as [ABA routing numbers](https://en.wikipedia.org/wiki/ABA_routing_transit_number), financial institution name lookup, and [Fedwire](https://en.wikipedia.org/wiki/Fedwire) and [FedACH](https://en.wikipedia.org/wiki/FedACH) routing information. The HTTP server is available in a [Docker image](#docker) and the Go package `github.com/moov-io/fed` is available. Moov's primary usage for this project is with ACH origination in our [paygate](https://github.com/moov-io/paygate) project.

The data and formats in this repository represent a compilation of **FedWire** and **FedACH** data from the [Federal Reserve Bank Services site](https://frbservices.org/). Both the official Fed plaintext and JSON file formats are supported.
