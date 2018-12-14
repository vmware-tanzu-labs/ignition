package http

import (
	_ "expvar" // metrics
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pivotalservices/ignition/api"
	"github.com/pivotalservices/ignition/config"
	"github.com/pivotalservices/ignition/http/session"
)

// API is the Ignition web app
type API struct {
	Ignition *config.Ignition
}

// URI is the combination of the scheme, domain, and port
func (a *API) URI() string {
	s := fmt.Sprintf("%s://%s", a.Ignition.Server.Scheme, a.Ignition.Server.Domain)
	if a.Ignition.Server.Port != 0 {
		s = fmt.Sprintf("%s:%v", s, a.Ignition.Server.Port)
	}
	return s
}

// Run starts a server listening on the given serveURI
func (a *API) Run() error {
	a.Ignition.Authorizer.Config.RedirectURL = fmt.Sprintf("%s%s", a.URI(), "/oauth2")
	r := a.createRouter()
	return http.ListenAndServe(fmt.Sprintf(":%v", a.Ignition.Server.ServePort), handlers.LoggingHandler(os.Stdout, handlers.CORS()(r)))
}

func (a *API) createRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(path.Join(a.Ignition.Server.WebRoot, "assets")+string(os.PathSeparator))))).Name("assets")
	r.Handle("/api/v1/profile", ensureHTTPClient(a.Ignition.Authorizer.SkipTLSValidation, ensureHTTPS(session.PopulateContext(Authenticate(api.ProfileHandler()), a.Ignition.Server.SessionStore))))
	infoHandler := api.InfoHandler(
		a.Ignition.Server.CompanyName,
		a.Ignition.Experimenter.SpaceName,
		a.Ignition.Experimenter.QuotaID,
		a.Ignition.Server.CollectAnalytics,
		a.Ignition.Experimenter.OrgCountUpdateInterval,
		a.Ignition.Deployment.CC)
	r.Handle("/api/v1/info", ensureHTTPClient(a.Ignition.Authorizer.SkipTLSValidation, Secure(infoHandler, a.Ignition.Authorizer.Domain, a.Ignition.Server.SessionStore)))

	orgHandler := api.OrganizationHandler(
		a.Ignition.Deployment.AppsURL,
		a.Ignition.Experimenter.OrgPrefix,
		a.Ignition.Experimenter.QuotaID,
		a.Ignition.Experimenter.ISOSegmentID,
		a.Ignition.Experimenter.SpaceName,
		a.Ignition.Deployment.CC)
	orgHandler = ensureUser(orgHandler, a.Ignition.Deployment.UAA, a.Ignition.Deployment.UAAOrigin, a.Ignition.Server.SessionStore)
	orgHandler = Secure(orgHandler, a.Ignition.Authorizer.Domain, a.Ignition.Server.SessionStore)
	r.Handle("/api/v1/organization", ensureHTTPClient(a.Ignition.Authorizer.SkipTLSValidation, orgHandler))

	a.handleAuth(r)
	r.Handle("/debug/vars", http.DefaultServeMux)
	var t *template.Template

	// If the API can't handle the route, let the SPA handle it
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if t == nil {
			var err error
			t, err = template.ParseFiles(filepath.Join(a.Ignition.Server.WebRoot, "index.html"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t.Execute(w, a.Ignition.Server)
	})

	return r
}
