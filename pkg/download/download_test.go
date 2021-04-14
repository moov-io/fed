// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package download

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestClient__fedach(t *testing.T) {
	client := setupClient(t)

	fedach, err := client.GetList("fedach")
	if err != nil {
		t.Fatal(err)
	}

	st, err := fedach.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if n := st.Size(); n < 1024 {
		t.Errorf("unexpected size of %d bytes", n)
	}

	bs, _ := ioutil.ReadAll(io.LimitReader(fedach, 10024))
	if !bytes.Contains(bs, []byte("fedACHParticipants")) {
		t.Errorf("unexpected output:\n%s", string(bs))
	}
}

func TestClient__fedwire(t *testing.T) {
	client := setupClient(t)

	fedwire, err := client.GetList("fedwire")
	if err != nil {
		t.Fatal(err)
	}

	st, err := fedwire.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if n := st.Size(); n < 1024 {
		t.Errorf("unexpected size of %d bytes", n)
	}

	bs, _ := ioutil.ReadAll(io.LimitReader(fedwire, 10024))
	if !bytes.Contains(bs, []byte("fedwireParticipants")) {
		t.Errorf("unexpected output:\n%s", string(bs))
	}
}

func setupClient(t *testing.T) *Client {
	t.Helper()

	routingNumber := os.Getenv("FRB_ROUTING_NUMBER")
	downloadCode := os.Getenv("FRB_DOWNLOAD_CODE")
	if routingNumber == "" || downloadCode == "" {
		t.Skip("missing FRB routing number or download code")
	}

	client, err := NewClient(&ClientOpts{
		RoutingNumber: routingNumber,
		DownloadCode:  downloadCode,
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}
