package main

import (
	"fmt"
	"github.com/moov-io/fed"
	"os"
)

// readFEDACHData opens and reads FedACHdir.txt then runs ACHDictionary.Read() to
// parse and define ACHDictionary properties
func (s *searcher) readFEDACHData() error {
	if s.logger != nil {
		s.logger.Log("read", "Read of FED data")
	}

	// ToDo: Path for main.go
	f, err := os.Open("../.././data/FedACHdir.txt")
	if err != nil {
		return fmt.Errorf("ERROR: opening FedACHdir.txt %v", err)
	}
	defer f.Close()

	s.ACHDictionary = fed.NewACHDictionary(f)
	if err := s.ACHDictionary.Read(); err != nil {
		return fmt.Errorf("ERROR: reading FedACHdir.txt %v", err)
	}

	if s.logger != nil {
		s.logger.Log("download", "Finished refresh of FED data")
	}

	return nil
}
