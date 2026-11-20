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

The HTTP server loads the plaintext (or official FRB JSON) directories in [`data/FedACHdir.txt`](https://github.com/moov-io/fed/tree/master/data/FedACHdir.txt) and [`data/fpddir.txt`](https://github.com/moov-io/fed/tree/master/data/fpddir.txt). The format of those files is documented in [`docs/FedACHdir.md`](https://github.com/moov-io/fed/tree/master/docs/FedACHdir.md) and [`docs/fpddir.md`](https://github.com/moov-io/fed/tree/master/docs/fpddir.md). The copies bundled in this repository are from 2018. The Federal Reserve Banks no longer publish the complete directories on the public website (as of December 9, 2018) and licensing prevents us from distributing more recent files. The Fed still compiles this data; see [Getting current directory data](https://github.com/moov-io/fed#getting-current-directory-data) in the project README.

Moov Fed can read the data files from anywhere on the filesystem. Set `FEDACH_DATA_PATH` / `FEDWIRE_DATA_PATH` or `INITIAL_DATA_DIRECTORY`. Both official formats from the Federal Reserve (plaintext and JSON) are supported. You can also run a [release binary](https://github.com/moov-io/fed/releases/latest) or build from source with `make build` (see [Binary distribution](usage-binary.md)).
