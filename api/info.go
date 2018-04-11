package api

import (
	"encoding/json"
	"net/http"
)

// Info is metadata that ignition API clients can use to display their UX
type Info struct {
	CompanyName              string
	ExperimentationSpaceName string
}

// InfoHandler writes the contents of the provided Info to the response
func InfoHandler(i Info) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(i)
	}
	return http.HandlerFunc(fn)
}
