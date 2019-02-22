package main

import (
	"fmt"
	"os"

	"github.com/moov-io/fed"
)

// readFEDACHData opens and reads FedACHdir.txt then runs ACHDictionary.Read() to
// parse and define ACHDictionary properties
func (s *searcher) readFEDACHData() error {
	if s.logger != nil {
		s.logger.Log("read", "Read of FED data")
	}

	f, err := os.Open("./data/FedACHdir.txt")
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

	f, err := os.Open("./data/fpddir.txt")
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
