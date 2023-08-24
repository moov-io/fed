// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/moov-io/base/log"
	"github.com/moov-io/fed"
	"github.com/moov-io/fed/pkg/download"
)

func fedACHDataFile(logger log.Logger) (io.Reader, error) {
	if file, err := attemptFileDownload(logger, "fedach"); file != nil {
		return file, nil
	} else if err != nil {
		return nil, fmt.Errorf("problem downloading fedach: %v", err)
	}

	path := readDataFilepath("FEDACH_DATA_PATH", "./data/FedACHdir.txt")
	logger.Logf("search: loading %s for ACH data", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("problem opening %s: %v", path, err)
	}
	return file, nil
}

func fedWireDataFile(logger log.Logger) (io.Reader, error) {
	if file, err := attemptFileDownload(logger, "fedwire"); file != nil {
		return file, nil
	} else if err != nil {
		return nil, fmt.Errorf("problem downloading fedwire: %v", err)
	}

	path := readDataFilepath("FEDWIRE_DATA_PATH", "./data/fpddir.txt")
	logger.Logf("search: loading %s for Wire data", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("problem opening %s: %v", path, err)
	}
	return file, nil
}

func attemptFileDownload(logger log.Logger, listName string) (io.Reader, error) {
	routingNumber := os.Getenv("FRB_ROUTING_NUMBER")
	downloadCode := os.Getenv("FRB_DOWNLOAD_CODE")

	if routingNumber != "" && downloadCode != "" {
		logger.Logf("download: attempting %s", listName)
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
func (s *searcher) readFEDACHData(reader io.Reader) error {
	if s.logger != nil {
		s.logger.Logf("Read of FED data")
	}

	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	s.ACHDictionary = fed.NewACHDictionary()
	if err := s.ACHDictionary.Read(reader); err != nil {
		return fmt.Errorf("ERROR: reading FedACHdir.txt %v", err)
	}

	if s.logger != nil {
		s.logger.Logf("Finished refresh of ACH data")
	}

	return nil
}

// readFEDWIREData opens and reads fpddir.txt then runs WIREDictionary.Read() to
// parse and define WIREDictionary properties
func (s *searcher) readFEDWIREData(reader io.Reader) error {
	if s.logger != nil {
		s.logger.Logf("Read of FED data")
	}

	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	s.WIREDictionary = fed.NewWIREDictionary()
	if err := s.WIREDictionary.Read(reader); err != nil {
		return fmt.Errorf("ERROR: reading fpddir.txt %v", err)
	}

	if s.logger != nil {
		s.logger.Logf("Finished refresh of WIRE data")
	}

	return nil
}
