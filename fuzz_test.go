package fed_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moov-io/fed"
)

func FuzzACHDictionary(f *testing.F) {
	populateCorpus(f, "FedACH", "fedach", ".txt", ".json")

	f.Fuzz(func(t *testing.T, contents string) {
		if len(contents) > 2<<20 {
			t.Skip()
		}

		dict := fed.NewACHDictionary()
		if err := dict.Read(strings.NewReader(contents)); err != nil {
			return
		}
		// Search APIs must not panic on loaded data
		_ = dict.RoutingNumberSearchSingle("021000021")
		_ = dict.FinancialInstitutionSearchSingle("JPMORGAN")
		_, _ = dict.RoutingNumberSearch("0210", 5)
		_ = dict.FinancialInstitutionSearch("BANK", 5)
		_ = dict.StateFilter("NY")
		_ = dict.CityFilter("NEW YORK")
	})
}

func FuzzWIREDictionary(f *testing.F) {
	populateCorpus(f, "fpddir", "FedWire", "wire", ".txt", ".json")

	f.Fuzz(func(t *testing.T, contents string) {
		if len(contents) > 2<<20 {
			t.Skip()
		}

		dict := fed.NewWIREDictionary()
		if err := dict.Read(strings.NewReader(contents)); err != nil {
			return
		}
		_ = dict.RoutingNumberSearchSingle("021000021")
		_ = dict.FinancialInstitutionSearchSingle("JPMORGAN")
		_, _ = dict.RoutingNumberSearch("0210", 5)
		_ = dict.FinancialInstitutionSearch("BANK", 5)
		_ = dict.StateFilter("NY")
		_ = dict.CityFilter("NEW YORK")
	})
}

func populateCorpus(f *testing.F, nameHints ...string) {
	f.Helper()

	f.Add("")
	f.Add("{}")
	f.Add("[]")

	_ = filepath.Walk("data", func(path string, info fs.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		base := strings.ToLower(filepath.Base(path))
		// Include all data files; nameHints just document intent
		_ = nameHints
		if strings.HasSuffix(base, ".txt") || strings.HasSuffix(base, ".json") {
			bs, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			// Fed directory files can be large; cap seeds
			if len(bs) > 512*1024 {
				// Still seed a prefix so format is represented
				f.Add(string(bs[:min(len(bs), 64*1024)]))
				return nil
			}
			f.Add(string(bs))
		}
		return nil
	})
}
