package main

import (
	"fmt"
	"os"

	"github.com/moov-io/fed"
)

var (
	fedACHDataFilepath  = os.Getenv("FEDACH_DATA_PATH")
	fedWIREDataFilepath = os.Getenv("FEDWIRE_DATA_PATH")
)

func init() {
	if fedACHDataFilepath == "" {
		fedACHDataFilepath = "./data/FedACHdir.txt"
	}
	if fedWIREDataFilepath == "" {
		fedWIREDataFilepath = "./data/fpddir.txt"
	}
}

// readFEDACHData opens and reads FedACHdir.txt then runs ACHDictionary.Read() to
// parse and define ACHDictionary properties
func (s *searcher) readFEDACHData() error {
	if s.logger != nil {
		s.logger.Log("read", "Read of FED data")
	}

	f, err := os.Open(fedACHDataFilepath)
	if err != nil {
		return fmt.Errorf("ERROR: opening FedACHdir.txt %v", err)
	}
	defer f.Close()

	s.ACHDictionary = fed.NewACHDictionary(f)
	if err := s.ACHDictionary.Read(); err != nil {
		return fmt.Errorf("ERROR: reading FedACHdir.txt %v", err)
	}

	if s.logger != nil {
		s.logger.Log("read", "Finished refresh of ACH data")
	}

	return nil
}

// readFEDWIREData opens and reads fpddir.txt then runs WIREDictionary.Read() to
// parse and define WIREDictionary properties
func (s *searcher) readFEDWIREData() error {
	if s.logger != nil {
		s.logger.Log("read", "Read of FED data")
	}

	f, err := os.Open(fedWIREDataFilepath)
	if err != nil {
		return fmt.Errorf("ERROR: opening fpddir.txt %v", err)
	}
	defer f.Close()

	s.WIREDictionary = fed.NewWIREDictionary(f)
	if err := s.WIREDictionary.Read(); err != nil {
		return fmt.Errorf("ERROR: reading fpddir.txt %v", err)
	}

	if s.logger != nil {
		s.logger.Log("read", "Finished refresh of WIRE data")
	}

	return nil
}
