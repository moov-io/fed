// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package logos

import (
	"os"
)

type Grabber interface {
	Lookup(name string) (*Logo, error)
}

type Logo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func NewGrabber() Grabber {
	clearbitApiKey := os.Getenv("CLEARBIT_API_KEY")
	if clearbitApiKey != "" {
		return newClearbit(clearbitApiKey)
	}
	return &noopGrabber{}
}

type noopGrabber struct{}

func NewNopGrabber() Grabber {
	return &noopGrabber{}
}

func (g *noopGrabber) Lookup(name string) (*Logo, error) {
	return nil, nil
}
