## Moov Fed

**[Documentation](https://moov-io.github.io/fed)** | **[Source](https://github.com/moov-io/fed)** | **[Running](https://github.com/moov-io/fed#usage)** | **[Configuration](https://github.com/moov-io/fed#configuration-settings)**

### Purpose

[Moov Fed](https://github.com/moov-io/fed) implements an HTTP interface to search [Fedwire](https://github.com/moov-io/fed/tree/master/docs/fpddir.md) and [FedACH](https://github.com/moov-io/fed/tree/master/docs/FedACHdir.md) data from the Federal Reserve Bank Services.

The data and formats below represent a compilation of  **Fedwire** and **FedACH** data from the [Federal Reserve Bank Services site](https://frbservices.org/):

* [FEDACH](https://github.com/moov-io/fed/tree/master/docs/FedACHdir.md)

* [FEDWire](https://github.com/moov-io/fed/tree/master/docs/fpddir.md)

Fed can be used standalone to search for routing numbers by Financial Institution name, city, state, postal code, and routing number. It can also be used in conjunction with [ACH](https://github.com/moov-io/ach) and [WIRE](https://github.com/moov-io/wire) to validate routing numbers.

## FED Data Files

The HTTP server loads the plaintext (or official FRB JSON) directories in [`data/FedACHdir.txt`](https://github.com/moov-io/fed/tree/master/data/FedACHdir.txt) and [`data/fpddir.txt`](https://github.com/moov-io/fed/tree/master/data/fpddir.txt). The format of those files is documented in [`FedACHdir.md`](FedACHdir.md) and [`fpddir.md`](fpddir.md). The copies bundled in this repository are from 2018. The Federal Reserve Banks no longer publish the complete directories on the public website (as of December 9, 2018) and licensing prevents us from distributing more recent files. The Fed still compiles this data; see [Getting current directory data](https://github.com/moov-io/fed#getting-current-directory-data) in the project README.

Moov Fed can read the data files from anywhere on the filesystem. Set `FEDACH_DATA_PATH` / `FEDWIRE_DATA_PATH` or `INITIAL_DATA_DIRECTORY`. Both official formats from the Federal Reserve (plaintext and JSON) are supported.

### Copyright and Terms of Use

Copyright &copy; Federal Reserve Banks

By accessing the [data](https://github.com/moov-io/fed/tree/master/data) in this repository you agree to the [Federal Reserve Banks' Terms of Use](https://frbservices.org/terms/index.html) and the [E-Payments Routing Directory Terms of Use Agreement](https://www.frbservices.org/EPaymentsDirectory/agreement.html).

## Disclaimer

**THIS REPOSITORY IS NOT AFFILIATED WITH THE FEDERAL RESERVE BANKS AND IS NOT AN OFFICIAL SOURCE FOR FEDWIRE AND FEDACH DATA.**

## Getting Help

 channel | info
 ------- | -------
 [Project Documentation](https://moov-io.github.io/fed/) | Our project documentation available online.
Twitter [@moov](https://twitter.com/moov)	| You can follow Moov.io's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io/fed/issues) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channel (`#fed`) to have an interactive discussion about the development of the project.
