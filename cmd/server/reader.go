// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/moov-io/base/log"
	"github.com/moov-io/fed"
	"github.com/moov-io/fed/pkg/download"
)

func fedACHDataFile(logger log.Logger) (io.Reader, error) {
	file, err := attemptFileDownload(logger, "fedach")
	if err != nil && !errors.Is(err, download.ErrMissingConfigValue) {
		return nil, fmt.Errorf("problem downloading fedach: %v", err)
	}

	if file != nil {
		logger.Info().Log("search: downloaded ACH file")

		return file, nil
	}

	path := readDataFilepath("FEDACH_DATA_PATH", "./data/FedACHdir.txt")
	logger.Logf("search: loading %s for ACH data", path)

	file, err = os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("problem opening %s: %v", path, err)
	}
	return file, nil
}

func fedWireDataFile(logger log.Logger) (io.Reader, error) {
	file, err := attemptFileDownload(logger, "fedwire")
	if err != nil && !errors.Is(err, download.ErrMissingConfigValue) {
		return nil, fmt.Errorf("problem downloading fedwire: %v", err)
	}

	if file != nil {
		logger.Info().Log("search: downloaded Wire file")

		return file, nil
	}

	path := readDataFilepath("FEDWIRE_DATA_PATH", "./data/fpddir.txt")
	logger.Logf("search: loading %s for Wire data", path)

	file, err = os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("problem opening %s: %v", path, err)
	}
	return file, nil
}

func attemptFileDownload(logger log.Logger, listName string) (io.Reader, error) {
	logger.Logf("download: attempting %s", listName)
	client, err := download.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("client setup: %w", err)
	}
	return client.GetList(listName)
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
		s.logger.Logf("Read of FED ACH data from %T", reader)
	}

	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	s.ACHDictionary = fed.NewACHDictionary()
	if err := s.ACHDictionary.Read(reader); err != nil {
		return fmt.Errorf("ERROR: reading FedACHdir.txt %v", err)
	}

	if s.logger != nil {
		s.logger.With(log.Fields{
			"records": log.Int(len(s.ACHDictionary.ACHParticipants)),
		}).Logf("Finished refresh of ACH data")
	}

	return nil
}

// readFEDWIREData opens and reads fpddir.txt then runs WIREDictionary.Read() to
// parse and define WIREDictionary properties
func (s *searcher) readFEDWIREData(reader io.Reader) error {
	if s.logger != nil {
		s.logger.Logf("Read of FED Wire data from %T", reader)
	}

	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	s.WIREDictionary = fed.NewWIREDictionary()
	if err := s.WIREDictionary.Read(reader); err != nil {
		return fmt.Errorf("ERROR: reading fpddir.txt %v", err)
	}

	if s.logger != nil {
		s.logger.With(log.Fields{
			"records": log.Int(len(s.WIREDictionary.WIREParticipants)),
		}).Logf("Finished refresh of WIRE data")
	}

	return nil
}
