package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Routing holds a FedACH routing record
type Routing struct {
	//RoutingNumber (9): 011000015
	RoutingNumber string
	// OfficeCode (1): O
	OfficeCode string
	// ServicingFrbNumber (9): 011000015
	ServicinFrbNumber string
	// RecordTypeCode (1): 0
	RecordTypeCode string
	// ChangeDate (6): 122415
	ChangeDate string
	// NewRoutingNumber (9): 000000000
	NewRoutingNumber string
	// CustomerName (36): FEDERAL RESERVE BANK
	CustomerName string
	// Address (36): 1000 PEACHTREE ST N.E.
	Address string
	// City (20): ATLANTA
	City string
	// State (2): GA
	State string
	// PostalCode (5): 30309
	PostalCode string
	// PostalCodeExtension (4): 4470
	PostalCodeExtension string
	// PhoneNumber(10): 8773722457
	PhoneNumber string
	// StatusCode (1): 1
	StatusCode string
	// ViewCode (1): 1
	ViewCode string
}

// Dictionary of Routing records accessable by Routes
type Dictionary struct {
	Routes []*Routing
}

func (d *Dictionary) addRouting(routing *Routing) *Dictionary {
	d.Routes = append(d.Routes, routing)
	return d
}

func readFile() {
	f, err := os.Open("./testdata/FedACHdir.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	dict := new(Dictionary)
	for scanner.Scan() {
		dict.addRouting(parseFedACH(scanner.Text()))
		fmt.Printf("Routing # %v \t Name: %v\n", dict.Routes[lineNum].RoutingNumber, dict.Routes[lineNum].CustomerName)
		lineNum++
	}
}

func parseFedACH(line string) *Routing {
	if len(line) == 0 {
		return nil
	}

	r := new(Routing)
	//RoutingNumber (9): 011000015
	r.RoutingNumber = line[:9]
	// OfficeCode (1): O
	r.OfficeCode = line[9:10]
	// ServicingFrbNumber (9): 011000015
	r.ServicinFrbNumber = line[10:19]
	// RecordTypeCode (1): 0
	r.RecordTypeCode = line[19:20]
	// ChangeDate (6): 122415
	r.ChangeDate = line[20:26]
	// NewRoutingNumber (9): 000000000
	r.NewRoutingNumber = line[26:35]
	// CustomerName (36): FEDERAL RESERVE BANK
	r.CustomerName = line[35:71]
	// Address (36): 1000 PEACHTREE ST N.E.
	r.Address = line[71:107]
	// City (20): ATLANTA
	r.City = line[107:127]
	// State (2): GA
	r.State = line[127:129]
	// PostalCode (5): 30309
	r.PostalCode = line[129:134]
	// PostalCodeExtension (4): 4470
	r.PostalCodeExtension = line[134:138]
	// PhoneNumber(10): 8773722457
	r.PhoneNumber = line[138:148]
	// StatusCode (1): 1
	r.StatusCode = line[148:149]
	// ViewCode (1): 1
	r.ViewCode = line[149:150]

	return r
}

func main() {
	readFile()
}
