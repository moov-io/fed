package feddir

import (
	"bufio"
	"io"
)

// RoutingDictionary of Participant records
type RoutingDictionary struct {
	Participants []*Participant
	scanner      *bufio.Scanner
	line         string
}

func NewRoutingDictionary(r io.Reader) *RoutingDictionary {
	return &RoutingDictionary{
		scanner: bufio.NewScanner(r),
	}
}

// Participant holds a FedACH routing record as defined by FEDACh
// https://www.frbservices.org/EPaymentsDirectory/achFormat.html
type Participant struct {
	// RoutingNumber The institution's routing number
	RoutingNumber string
	// OfficeCode Main office or branch O=main B=branch
	OfficeCode string
	// ServicingFrbNumber Servicing Fed's main office routing number
	ServicinFrbNumber string
	// RecordTypeCode The code indicating the ABA number to be used to route or send ACH items to the RFI
	// 0 = Institution is a Federal Reserve Bank
	// 1 = Send items to customer routing number
	// 2 = Send items to customer using new routing number field
	RecordTypeCode string
	// Revised Date of last revision: YYYYMMDD, or blank
	Revised string
	// NewRoutingNumber Institution's new routing number resulting from a merger or renumber
	NewRoutingNumber string
	// CustomerName (36): FEDERAL RESERVE BANK
	CustomerName string
	// Location is the delivery address
	Location
	// PhoneNumber The institution's phone number
	PhoneNumber string
	// StatusCode Code is based on the customers receiver code
	// 1=Receives Gov/Comm
	StatusCode string
	// ViewCode (1): 1
	ViewCode string
}

// Location City name and state code in the institution's delivery address
type Location struct {
	// Address
	Address string
	// City
	City string
	// State
	State string
	// PostalCode
	PostalCode string
	// PostalCodeExtension
	PostalCodeExtension string
}

// Read parses a single line or multiple line
func (rd *RoutingDictionary) Read() error {
	// read through the entire file
	for rd.scanner.Scan() {
		rd.line = rd.scanner.Text()

		if err := rd.parseRouting(); err != nil {
			return err
		}
	}
	return nil
}

func (rd *RoutingDictionary) parseRouting() error {
	route := new(Participant)

	if len(rd.line) == 0 {
		return nil
	}

	//RoutingNumber (9): 011000015
	route.RoutingNumber = rd.line[:9]
	// OfficeCode (1): O
	route.OfficeCode = rd.line[9:10]
	// ServicingFrbNumber (9): 011000015
	route.ServicinFrbNumber = rd.line[10:19]
	// RecordTypeCode (1): 0
	route.RecordTypeCode = rd.line[19:20]
	// ChangeDate (6): 122415
	route.Revised = rd.line[20:26]
	// NewRoutingNumber (9): 000000000
	route.NewRoutingNumber = rd.line[26:35]
	// CustomerName (36): FEDERAL RESERVE BANK
	route.CustomerName = rd.line[35:71]
	// Address (36): 1000 PEACHTREE ST N.E.
	route.Address = rd.line[71:107]
	// City (20): ATLANTA
	route.City = rd.line[107:127]
	// State (2): GA
	route.State = rd.line[127:129]
	// PostalCode (5): 30309
	route.PostalCode = rd.line[129:134]
	// PostalCodeExtension (4): 4470
	route.PostalCodeExtension = rd.line[134:138]
	// PhoneNumber(10): 8773722457
	route.PhoneNumber = rd.line[138:148]
	// StatusCode (1): 1
	route.StatusCode = rd.line[148:149]
	// ViewCode (1): 1
	route.ViewCode = rd.line[149:150]

	rd.Participants = append(rd.Participants, route)

	return nil

}
