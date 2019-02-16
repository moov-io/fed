// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package fed

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	fedFilenames = []string{
		"FedACHdir.txt",          // FEDACH
		"fpddir.txt",          // FEDWIRE
	}
	// ToDo:  Find out where to get the files
	fedURLTemplate = ""
)

func init() {
	if v := os.Getenv("FED_DOWNLOAD_TEMPLATE"); v != "" {
		fedURLTemplate = v
	}
}

// Downloader will download and cache fed files in a temp directory.
//
// If HTTP is nil then http.DefaultClient will be used (which has NO timeouts).
type Downloader struct {
	HTTP *http.Client
}

// GetFiles will download all FED related files and store them in a temporary directory
// returned and an error otherwise.
//
// Callers are expected to cleanup the temp directory.
func (dl *Downloader) GetFiles() (string, error) {
	if dl.HTTP == nil {
		dl.HTTP = http.DefaultClient
	}

	dir, err := ioutil.TempDir("", "fed-downloader")
	if err != nil {
		return "", fmt.Errorf("FED: unable to make temp dir: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(fedFilenames))

	for i := range fedFilenames {
		name := fedFilenames[i]

		go func(wg *sync.WaitGroup, filename string) {
			defer wg.Done()

			resp, err := dl.HTTP.Get(fmt.Sprintf(fedURLTemplate, filename))
			if err != nil {
				return
			}
			defer resp.Body.Close()

			// Copy resp.Body into a file in our temp dir
			fd, err := os.Create(filepath.Join(dir, filename))
			if err != nil {
				return
			}
			defer fd.Close()

			io.Copy(fd, resp.Body) // copy contents
		}(&wg, name)
	}

	wg.Wait()

	// count files and error if the count isn't what we expected
	fds, err := ioutil.ReadDir(dir)
	if err != nil || len(fds) != len(fedFilenames) {
		return "", fmt.Errorf("FED: problem downloading (found=%d, expected=%d): err=%v", len(fds), len(fedFilenames), err)
	}

	return dir, nil
}

