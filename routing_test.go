package feddir

import (
	"strings"
	"testing"
)

func TestParseRouting(t *testing.T) {
	var line = "073905527O0710003011012908000000000LINCOLN SAVINGS BANK                P O BOX E                           REINBECK            IA506690159319788644111     "
	//fmt.Println("wade got here: " + line)

	rd := NewRoutingDictionary(strings.NewReader(line))
	rd.Read()
	route := rd.Participants[0]
	if route.CustomerName != "LINCOLN SAVINGS BANK" {
		t.Errorf("CustomerName Expected 'LINCOLN SAVINGS BANK' got: %v", route.CustomerName)
	}
	/*
		r := NewReader(strings.NewReader(line))
		r.line = line
		if err := r.parseEntryDetail(); err != nil {
			t.Errorf("%T: %s", err, err)
		}
		record := r.currentBatch.GetEntries()[0]

		if record.recordType != "6" {
			t.Errorf("RecordType Expected '6' got: %v", record.recordType)
		}
		if record.TransactionCode != 27 {
			t.Errorf("TransactionCode Expected '27' got: %v", record.TransactionCode)
		}
		if record.RDFIIdentificationField() != "05320001" {
			t.Errorf("RDFIIdentification Expected '05320001' got: '%v'", record.RDFIIdentificationField())
		}
		if record.CheckDigit != "9" {
			t.Errorf("CheckDigit Expected '9' got: %v", record.CheckDigit)
		}
		if record.DFIAccountNumberField() != "12345            " {
			t.Errorf("DfiAccountNumber Expected '12345            ' got: %v", record.DFIAccountNumberField())
		}
		if record.AmountField() != "0000010500" {
			t.Errorf("Amount Expected '0000010500' got: %v", record.AmountField())
		}

		if record.IdentificationNumber != "c-1            " {
			t.Errorf("IdentificationNumber Expected 'c-1            ' got: %v", record.IdentificationNumber)
		}
		if record.IndividualName != "Arnold Wade           " {
			t.Errorf("IndividualName Expected 'Arnold Wade           ' got: %v", record.IndividualName)
		}
		if record.DiscretionaryData != "DD" {
			t.Errorf("DiscretionaryData Expected 'DD' got: %v", record.DiscretionaryData)
		}
		if record.AddendaRecordIndicator != 0 {
			t.Errorf("AddendaRecordIndicator Expected '0' got: %v", record.AddendaRecordIndicator)
		}
		if record.TraceNumberField() != "076401255655291" {
			t.Errorf("TraceNumber Expected '076401255655291' got: %v", record.TraceNumberField())
		}
	*/
}

/*
func testPPDDebitRead(t testing.TB) {
	f, err := os.Open("./data/testdata/ppd-debit.ach")
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}
	defer f.Close()
	r := NewReader(f)
	_, err = r.Read()
	if err != nil {
		t.Errorf("%T: %s", err, err)
	}
	if err = r.File.Validate(); err != nil {
		t.Errorf("%T: %s", err, err)
	}
}
*/

/*

func readFile() map[string]*Routing {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Printf("reading file: %v", f.Name())
	scanner := bufio.NewScanner(f)
	lineNum := 0
	lookup = make(map[string]*Routing)
	//dict := new(Dictionary)
	for scanner.Scan() {
		//dict.addRouting(parseRouting(scanner.Text()))
		route := parseRouting(scanner.Text())
		lookup[route.RoutingNumber] = route
		//fmt.Printf("Routing # %v \t Name: %v\n", dict.Routes[lineNum].RoutingNumber, dict.Routes[lineNum].CustomerName)
		lineNum++
	}
	return lookup
}
*/
