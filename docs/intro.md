---
layout: page
title: Intro
hide_hero: true
show_sidebar: false
menubar: docs-menu
---

## What is Moov Fed?

[Moov Fed](https://github.com/moov-io/fed) implements an HTTP interface to search [Fedwire](https://github.com/moov-io/fed/tree/master/docs/fpddir.md) and [FedACH](https://github.com/moov-io/fed/tree/master/docs/FedACHdir.md) data from the Federal Reserve Bank Services.

The data and formats below represent a compilation of  **Fedwire** and **FedACH** data from the [Federal Reserve Bank Services site](https://frbservices.org/):

* [FEDACH](https://github.com/moov-io/fed/tree/master/docs/FedACHdir.md)

* [FEDWire](https://github.com/moov-io/fed/tree/master/docs/fpddir.md)

Fed can be used standalone to search for routing numbers by Financial Institution name, city, state, postal code, and routing number. It can also be used in conjunction with [Moov ACH](https://github.com/moov-io/ach) and [Moov Wire](https://github.com/moov-io/wire) to validate routing numbers.

## Data files

The data files included in this repository ([`FedACHdir.md`](https://github.com/moov-io/fed/tree/master/docs/FedACHdir.md) and [`fpddir.md`](https://github.com/moov-io/fed/tree/master/docs/fpddir.md)) are **outdated** and from 2018. The Fed no longer releases this data publicly and licensing on more recent files prevents us from distributing them. However, the Fed still complies this data and you can retrieve up-to-date files for use in our project, either from [LexisNexis](https://risk.lexisnexis.com/financial-services/payments-efficiency/payment-routing) or your financial institution.

Moov Fed can read the data files from anywhere on the filesystem. This allows you to mount the files and set `FEDACH_DATA_PATH` / `FEDWIRE_DATA_PATH` environmental variables. Both official formats from the Federal Reserve (plaintext and JSON) are supported.
