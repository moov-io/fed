// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package download

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClient__fedach(t *testing.T) {
	client := setupClient(t)

	fedach, err := client.GetList("fedach")
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := io.ReadAll(io.LimitReader(fedach, 10024))
	if !bytes.Contains(bs, []byte("fedACHParticipants")) {
		t.Errorf("unexpected output:\n%s", string(bs))
	}
}

func TestClient__fedach_custom_url(t *testing.T) {
	file, err := os.ReadFile(filepath.Join("..", "..", "data", "fedachdir.json"))
	if err != nil {
		t.Fatal(err)
	}

	mockHTTPServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, string(file))
	}))
	defer mockHTTPServer.Close()

	t.Setenv("FRB_DOWNLOAD_URL_TEMPLATE", mockHTTPServer.URL+"/%s")
	t.Setenv("FRB_ROUTING_NUMBER", "123456789")
	t.Setenv("FRB_DOWNLOAD_CODE", "a1b2c3d4-123b-9876-1234-z1x2y3a1b2c3")

	client := setupClient(t)

	fedach, err := client.GetList("fedach")
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := io.ReadAll(io.LimitReader(fedach, 10024))
	if !bytes.Equal(bs, file) {
		t.Errorf("unexpected output:\n%s", string(bs))
	}
}

func TestClient__fedwire(t *testing.T) {
	client := setupClient(t)

	fedwire, err := client.GetList("fedwire")
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := io.ReadAll(io.LimitReader(fedwire, 10024))
	if !bytes.Contains(bs, []byte("fedwireParticipants")) {
		t.Errorf("unexpected output:\n%s", string(bs))
	}
}

func TestClient__wire_custom_url(t *testing.T) {
	file, err := os.ReadFile(filepath.Join("..", "..", "data", "fpddir.json"))
	if err != nil {
		t.Fatal(err)
	}
	mockHTTPServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/fedwire" {
			writer.WriteHeader(http.StatusNotFound)
		}
		fmt.Fprint(writer, string(file))
	}))
	defer mockHTTPServer.Close()

	t.Setenv("FRB_DOWNLOAD_URL_TEMPLATE", mockHTTPServer.URL+"/%s")
	t.Setenv("FRB_ROUTING_NUMBER", "123456789")
	t.Setenv("FRB_DOWNLOAD_CODE", "a1b2c3d4-123b-9876-1234-z1x2y3a1b2c3")

	client := setupClient(t)

	fedwire, err := client.GetList("fedwire")
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := io.ReadAll(io.LimitReader(fedwire, 10024))
	if !bytes.Equal(bs, file) {
		t.Errorf("unexpected output:\n%s", string(bs))
	}
}

func setupClient(t *testing.T) *Client {
	t.Helper()

	routingNumber := os.Getenv("FRB_ROUTING_NUMBER")
	downloadCode := os.Getenv("FRB_DOWNLOAD_CODE")
	downloadURL := os.Getenv("FRB_DOWNLOAD_URL_TEMPLATE")
	if routingNumber == "" || downloadCode == "" {
		t.Skip("missing FRB routing number or download code")
	}

	client, err := NewClient(&ClientOpts{
		RoutingNumber: routingNumber,
		DownloadCode:  downloadCode,
		DownloadURL:   downloadURL,
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestNewClient__opts(t *testing.T) {
	t.Setenv("FRB_ROUTING_NUMBER", "")
	t.Setenv("FRB_DOWNLOAD_CODE", "")
	t.Setenv("FRB_DOWNLOAD_URL_TEMPLATE", "")

	var gotOrgID, gotDownloadCode string
	mockHTTPServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/fedach" {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		gotOrgID = request.Header.Get("X_FRB_EPAYMENTS_DIRECTORY_ORG_ID")
		gotDownloadCode = request.Header.Get("X_FRB_EPAYMENTS_DIRECTORY_DOWNLOAD_CD")
		fmt.Fprint(writer, "fedACHParticipants")
	}))
	defer mockHTTPServer.Close()

	client, err := NewClient(&ClientOpts{
		RoutingNumber: "123456789",
		DownloadCode:  "a1b2c3d4-123b-9876-1234-z1x2y3a1b2c3",
		DownloadURL:   mockHTTPServer.URL + "/%s",
	})
	if err != nil {
		t.Fatal(err)
	}

	fedach, err := client.GetList("fedach")
	if err != nil {
		t.Fatal(err)
	}
	bs, _ := io.ReadAll(fedach)
	if string(bs) != "fedACHParticipants" {
		t.Errorf("unexpected output: %q", string(bs))
	}
	if gotOrgID != "123456789" {
		t.Errorf("unexpected org id header: %q", gotOrgID)
	}
	if gotDownloadCode != "a1b2c3d4-123b-9876-1234-z1x2y3a1b2c3" {
		t.Errorf("unexpected download code header: %q", gotDownloadCode)
	}
}

func TestNewClient__optsOverrideEnv(t *testing.T) {
	t.Setenv("FRB_ROUTING_NUMBER", "987654321")
	t.Setenv("FRB_DOWNLOAD_CODE", "env-code")
	t.Setenv("FRB_DOWNLOAD_URL_TEMPLATE", "https://env.example.com/%s")

	client, err := NewClient(&ClientOpts{
		RoutingNumber: "123456789",
		DownloadCode:  "opts-code",
		DownloadURL:   "https://opts.example.com/%s",
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.routingNumber != "123456789" || client.downloadCode != "opts-code" || client.downloadURL != "https://opts.example.com/%s" {
		t.Errorf("unexpected client config: %#v", client)
	}
}

func TestNewClient__envOnly(t *testing.T) {
	t.Setenv("FRB_ROUTING_NUMBER", "123456789")
	t.Setenv("FRB_DOWNLOAD_CODE", "env-code")
	t.Setenv("FRB_DOWNLOAD_URL_TEMPLATE", "https://env.example.com/%s")

	tests := []struct {
		name string
		opts *ClientOpts
	}{
		{name: "nil", opts: nil},
		{name: "empty", opts: &ClientOpts{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.opts)
			if err != nil {
				t.Fatal(err)
			}
			if client.routingNumber != "123456789" || client.downloadCode != "env-code" || client.downloadURL != "https://env.example.com/%s" {
				t.Errorf("unexpected client config: %#v", client)
			}
		})
	}
}

func TestNewClient__missingConfig(t *testing.T) {
	t.Setenv("FRB_ROUTING_NUMBER", "")
	t.Setenv("FRB_DOWNLOAD_CODE", "")
	t.Setenv("FRB_DOWNLOAD_URL_TEMPLATE", "")

	if _, err := NewClient(nil); !errors.Is(err, ErrMissingRoutingNumber) {
		t.Errorf("expected ErrMissingRoutingNumber, got %v", err)
	}
	if _, err := NewClient(&ClientOpts{RoutingNumber: "123456789"}); !errors.Is(err, ErrMissingDownloadCD) {
		t.Errorf("expected ErrMissingDownloadCD, got %v", err)
	}
}

func TestClient_GetListRejectedDownloadCode(t *testing.T) {
	// FRB returns HTTP 200 with application-level code 202 and an empty list
	// when the routing number or download code is missing or incorrect.
	body := `{
  "fedACHParticipants" : {
    "response" : {
      "code" : 202
    },
    "fedACHParticipants" : [ ]
  }
}`
	var logs []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, body)
	}))
	defer srv.Close()

	t.Setenv("FRB_DOWNLOAD_URL_TEMPLATE", srv.URL+"/%s")
	t.Setenv("FRB_ROUTING_NUMBER", "123456780")
	t.Setenv("FRB_DOWNLOAD_CODE", "invalid-download-code")

	client, err := NewClient(&ClientOpts{
		Logf: func(format string, args ...any) {
			logs = append(logs, fmt.Sprintf(format, args...))
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetList("fedach")
	if err == nil {
		t.Fatal("expected error for rejected download code")
	}
	if !errors.Is(err, ErrEmptyDirectory) {
		t.Fatalf("expected ErrEmptyDirectory, got %v", err)
	}
	if !strings.Contains(err.Error(), "response code 202") {
		t.Fatalf("expected response code in error, got %v", err)
	}

	joined := strings.Join(logs, "\n")
	if !strings.Contains(joined, "HTTP 200") {
		t.Fatalf("expected HTTP status to be logged, got %q", joined)
	}
	if !strings.Contains(joined, "response code 202") {
		t.Fatalf("expected FRB response code to be logged, got %q", joined)
	}
}

func TestClient_GetListRejectedWireDownloadCode(t *testing.T) {
	body := `{
  "fedwireParticipants" : {
    "response" : {
      "code" : 202
    },
    "fedwireParticipants" : [ ]
  }
}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, body)
	}))
	defer srv.Close()

	t.Setenv("FRB_DOWNLOAD_URL_TEMPLATE", srv.URL+"/%s")
	t.Setenv("FRB_ROUTING_NUMBER", "123456780")
	t.Setenv("FRB_DOWNLOAD_CODE", "invalid-download-code")

	client, err := NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.GetList("fedwire")
	if err == nil {
		t.Fatal("expected error for rejected download code")
	}
	if !strings.Contains(err.Error(), "response code 202") {
		t.Fatalf("expected response code in error, got %v", err)
	}
}

func TestInspectFRBListSuccess(t *testing.T) {
	file, err := os.ReadFile(filepath.Join("..", "..", "data", "fedachdir.json"))
	if err != nil {
		t.Fatal(err)
	}
	code, count, err := inspectFRBList("fedach", file)
	if err != nil {
		t.Fatal(err)
	}
	if code != FRBSuccessResponseCode {
		t.Errorf("code: got %d", code)
	}
	if count == 0 {
		t.Error("expected participants")
	}
}

func TestInspectFRBListPlaintext(t *testing.T) {
	code, count, err := inspectFRBList("fedach", []byte("011000015O0110000150"))
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 || count != 0 {
		t.Errorf("plaintext should be ignored, got code=%d count=%d", code, count)
	}
}
