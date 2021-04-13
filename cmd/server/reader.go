// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/moov-io/fed"
	"github.com/moov-io/fed/pkg/download"

	"github.com/go-kit/kit/log"
)

func fedACHDataFile(logger log.Logger) *os.File {
	if file, err := attemptFileDownload(logger, "fedach"); file != nil {
		return file
	} else if err != nil {
		panic(fmt.Sprintf("problem downloading fedach: %v", err))
	}

	path := readDataFilepath("FEDACH_DATA_PATH", "./data/FedACHdir.txt")
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("problem opening %s: %v", path, err))
	}
	return file
}

func fedWireDataFile(logger log.Logger) *os.File {
	if file, err := attemptFileDownload(logger, "fedwire"); file != nil {
		return file
	} else if err != nil {
		panic(fmt.Sprintf("problem downloading fedwire: %v", err))
	}

	path := readDataFilepath("FEDWIRE_DATA_PATH", "./data/fpddir.txt")
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("problem opening %s: %v", path, err))
	}
	return file
}

func attemptFileDownload(logger log.Logger, listName string) (*os.File, error) {
	routingNumber := os.Getenv("FRB_ROUTING_NUMBER")
	downloadCode := os.Getenv("FRB_DOWNLOAD_CODE")

	if routingNumber != "" && downloadCode != "" {
		logger.Log("download", fmt.Sprintf("attempting %s download", listName))
		client, err := download.NewClient(&download.ClientOpts{
			RoutingNumber: routingNumber,
			DownloadCode:  downloadCode,
		})
		if err != nil {
			return nil, fmt.Errorf("client setup: %v", err)
		}
		return client.GetList(listName)
	}

	return nil, nil
}

func readDataFilepath(env, fallback string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return fallback
}

// readFEDACHData opens and reads FedACHdir.txt then runs ACHDictionary.Read() to
// parse and define ACHDictionary properties
func (s *searcher) readFEDACHData(file *os.File) error {
	if s.logger != nil {
		s.logger.Logf("Read of FED data")
	}
	defer file.Close()

	s.ACHDictionary = fed.NewACHDictionary()
	if err := s.ACHDictionary.Read(file); err != nil {
		return fmt.Errorf("ERROR: reading FedACHdir.txt %v", err)
	}

	if s.logger != nil {
		s.logger.Logf("Finished refresh of ACH data")
	}

	return nil
}

// readFEDWIREData opens and reads fpddir.txt then runs WIREDictionary.Read() to
// parse and define WIREDictionary properties
func (s *searcher) readFEDWIREData(file *os.File) error {
	if s.logger != nil {
		s.logger.Logf("Read of FED data")
	}
	defer file.Close()

	s.WIREDictionary = fed.NewWIREDictionary()
	if err := s.WIREDictionary.Read(file); err != nil {
		return fmt.Errorf("ERROR: reading fpddir.txt %v", err)
	}

	if s.logger != nil {
		s.logger.Logf("Finished refresh of WIRE data")
	}

	return nil
}
