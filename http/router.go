package http

import (
	_ "expvar" // metrics
	"fmt"
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
	r.Handle("/api/v1/profile", ensureHTTPS(session.PopulateContext(Authorize(api.ProfileHandler(), a.Ignition.Authorizer.Domain), a.Ignition.Server.SessionStore)))
	r.Handle("/api/v1/info", ensureHTTPS(api.InfoHandler(api.Info{
		CompanyName:              a.Ignition.Server.CompanyName,
		ExperimentationSpaceName: a.Ignition.Experimenter.SpaceName,
	})))

	orgHandler := api.OrganizationHandler(a.Ignition.Deployment.AppsURL, a.Ignition.Experimenter.OrgPrefix, a.Ignition.Experimenter.QuotaID, a.Ignition.Experimenter.SpaceName, a.Ignition.Deployment.CC)
	orgHandler = ensureUser(orgHandler, a.Ignition.Deployment.UAA, a.Ignition.Deployment.UAAOrigin, a.Ignition.Server.SessionStore)
	orgHandler = Authorize(orgHandler, a.Ignition.Authorizer.Domain)
	orgHandler = session.PopulateContext(orgHandler, a.Ignition.Server.SessionStore)
	orgHandler = ensureHTTPS(orgHandler)
	r.Handle("/api/v1/organization", orgHandler)

	a.handleAuth(r)
	r.Handle("/debug/vars", http.DefaultServeMux)

	// If the API can't handle the route, let the SPA handle it
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, filepath.Join(a.Ignition.Server.WebRoot, "index.html"))
	})

	return r
}
