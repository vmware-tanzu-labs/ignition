package internal

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// HandleTestdata is a http.Handler that loads the given filename from the
// testdata directory and serves it over HTTP.
func HandleTestdata(t *testing.T, s string, called func()) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		f, err := os.Open(filepath.Join("testdata", s))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		io.Copy(w, f)
	})
}

// ServeFromTestdata is a test HTTP server that serves files from the testdata directory.
func ServeFromTestdata(t *testing.T, s string, called func()) *httptest.Server {
	t.Helper()
	return httptest.NewServer(HandleTestdata(t, s, called))
}

// StringFromTestdata reads the given file from the testdata dierctory and
// returns it as a string
func StringFromTestdata(t *testing.T, s string) string {
	f, err := os.Open(filepath.Join("testdata", s))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
