package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pivotalservices/ignition/user"
)

func ProfileHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		p, err := user.ProfileFromContext(req.Context())
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
	return http.HandlerFunc(fn)
}
