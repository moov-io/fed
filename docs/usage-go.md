---
layout: page
title: Go library
hide_hero: true
show_sidebar: false
menubar: docs-menu
---

# Go library

[![GoDoc](https://godoc.org/github.com/moov-io/fed?status.svg)](https://godoc.org/github.com/moov-io/fed) [![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/fed)](https://goreportcard.com/report/github.com/moov-io/fed)

Fed can be included in your application or run as a standalone HTTP server.

## Go module

Latest stable Release for Go: [Project Releases](https://github.com/moov-io/fed/releases)

To get started with the Fed Go library, first import the module in your Go source:

```go
import (
    "github.com/moov-io/fed"
)
```

And then install or update the project with:

```
$ go get -u github.com/moov-io/fed
```

Requires [Go v1.26 or newer](https://golang.org/doc/install).

## Search the FedACH directory

This example loads the bundled plaintext FedACH file and looks up a routing number. Run it from a checkout (or change the path) so `data/FedACHdir.txt` is reachable.

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

`ACHDictionary.RoutingNumberSearch(prefix, limit)` searches by routing-number prefix (at least the first two digits). `WIREDictionary` has the same methods for Fedwire participants.

## HTTP client

The generated client in [`github.com/moov-io/fed/client`](https://pkg.go.dev/github.com/moov-io/fed/client) talks to a running Fed HTTP server.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/antihax/optional"
	"github.com/moov-io/fed/client"
)

func main() {
	cfg := client.NewConfiguration()
	cfg.BasePath = "http://localhost:8086"
	api := client.NewAPIClient(cfg)

	resp, _, err := api.FEDApi.SearchFEDACH(context.Background(), &client.SearchFEDACHOpts{
		RoutingNumber: optional.NewString("273976369"),
	})
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.ACHParticipants) == 0 {
		log.Fatal("expected a participant")
	}
	p := resp.ACHParticipants[0]
	fmt.Printf("%s %s\n", p.RoutingNumber, p.CustomerName)
	// 273976369 VERIDIAN CREDIT UNION
}
```

Start the server first (`make build && ./bin/server`, or a release binary with data files). See [Binary distribution](usage-binary.md).
