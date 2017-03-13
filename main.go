package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
)

const (
	batchSize = 100
	filePath  = "./data/FedACHdir.txt"
	indexPath = "routing-search.bleve"
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

func readFile() *Dictionary {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Printf("reading file: %v", f.Name())
	scanner := bufio.NewScanner(f)
	lineNum := 0
	dict := new(Dictionary)
	for scanner.Scan() {
		dict.addRouting(parseRouting(scanner.Text()))
		//fmt.Printf("Routing # %v \t Name: %v\n", dict.Routes[lineNum].RoutingNumber, dict.Routes[lineNum].CustomerName)
		lineNum++
	}
	return dict
}

func parseRouting(line string) *Routing {
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
	routingIndex, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Printf("Creating new index...")
		// create a mapping
		indexMapping, err := buildIndexMapping()
		if err != nil {
			fmt.Println("buildIndexMapping")
			log.Fatal(err)
		}
		//indexMapping := bleve.NewIndexMapping()
		routingIndex, err = bleve.New(indexPath, indexMapping)
		if err != nil {
			fmt.Println("routingIndex")
			log.Fatal(err)
		}
		log.Println("go routine")
		err = indexRouting(routingIndex)
		if err != nil {
			log.Println("go routine fatal")

			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Opening existing index...\n")
	}
	query := bleve.NewMatchQuery("veridian")
	search := bleve.NewSearchRequest(query)
	searchResults, err := routingIndex.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)
}

func indexRouting(i bleve.Index) error {
	// build the Dictionary
	log.Printf("Indexing...\n")
	dict := readFile()

	count := 0
	batch := i.NewBatch()
	batchCount := 0

	for _, route := range dict.Routes {
		batch.Index(route.RoutingNumber, route)
		batchCount++
		if batchCount >= batchSize {
			err := i.Batch(batch)
			if err != nil {
				return err
			}
			batch = i.NewBatch()
			batchCount = 0
			fmt.Printf("Created batch index at count: %v \n", count)
		}
		count++

	}
	if batchCount > 0 {
		err := i.Batch(batch)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func buildIndexMapping() (mapping.IndexMapping, error) {

	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	//keywordFieldMapping := bleve.NewTextFieldMapping()
	//keywordFieldMapping.Analyzer = keyword.Name

	routingMapping := bleve.NewDocumentMapping()

	// name
	routingMapping.AddFieldMappingsAt("CustomerName", englishTextFieldMapping)

	//routingMapping.AddFieldMappingsAt("City", keywordFieldMapping)
	//routingMapping.AddFieldMappingsAt("State", keywordFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("Routing", routingMapping)

	//indexMapping.TypeField = "type"
	//indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}
