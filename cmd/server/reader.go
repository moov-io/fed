// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/moov-io/fed"
)

var (
	fedACHDataFilepath = func() string {
		return readDataFilepath("FEDACH_DATA_PATH", "./data/FedACHdir.txt")
	}()
	fedWIREDataFilepath = func() string {
		return readDataFilepath("FEDWIRE_DATA_PATH", "./data/fpddir.txt")
	}()
)

func readDataFilepath(env, fallback string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return fallback
}

// readFEDACHData opens and reads FedACHdir.txt then runs ACHDictionary.Read() to
// parse and define ACHDictionary properties
func (s *searcher) readFEDACHData(path string) error {
	if s.logger != nil {
		s.logger.Logf("Read of FED data")
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("ERROR: opening FedACHdir.txt %v", err)
	}
	defer f.Close()

	s.ACHDictionary = fed.NewACHDictionary()
	if err := s.ACHDictionary.Read(f); err != nil {
		return fmt.Errorf("ERROR: reading FedACHdir.txt %v", err)
	}

	if s.logger != nil {
		s.logger.Logf("Finished refresh of ACH data")
	}

	return nil
}

// readFEDWIREData opens and reads fpddir.txt then runs WIREDictionary.Read() to
// parse and define WIREDictionary properties
func (s *searcher) readFEDWIREData(path string) error {
	if s.logger != nil {
		s.logger.Logf("Read of FED data")
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("ERROR: opening fpddir.txt %v", err)
	}
	defer f.Close()

	s.WIREDictionary = fed.NewWIREDictionary()
	if err := s.WIREDictionary.Read(f); err != nil {
		return fmt.Errorf("ERROR: reading fpddir.txt %v", err)
	}

	if s.logger != nil {
		s.logger.Logf("Finished refresh of WIRE data")
	}

	return nil
}
