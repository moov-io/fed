// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package logos

type Grabber interface {
	Lookup(name string) (*Logo, error)
}

type Logo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
