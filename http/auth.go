package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/gologin"
	dgoauth2 "github.com/dghubble/gologin/oauth2"
	"github.com/dghubble/sessions"
	"github.com/gorilla/mux"
	"github.com/pivotalservices/ignition/http/session"
	"github.com/pivotalservices/ignition/uaa"
	"github.com/pivotalservices/ignition/user"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func (a *API) handleAuth(r *mux.Router) {
	stateConfig := gologin.DefaultCookieConfig
	if a.Ignition.Server.Domain == "localhost" {
		stateConfig = gologin.DebugOnlyCookieConfig
	}
	r.Handle("/login", ensureHTTPS(dgoauth2.StateHandler(stateConfig, dgoauth2.LoginHandler(a.Ignition.Authorizer.Config, nil)))).Name("login")
	r.Handle("/oauth2", ensureHTTPS(dgoauth2.StateHandler(stateConfig, CallbackHandler(a.Ignition.Authorizer.Config, a.Ignition.Authorizer.Fetcher, session.IssueSession(a.Ignition.Server.SessionStore, a.Ignition.Deployment.UAA), session.LogoutHandler(a.Ignition.Server.SessionStore))))).Name("oauth2")
	r.Handle("/logout", ensureHTTPS(session.LogoutHandler(a.Ignition.Server.SessionStore))).Name("logout")
}

func ensureUser(next http.Handler, uaa uaa.API, origin string, s sessions.Store) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID, err := session.UserIDFromContext(r.Context())
		if strings.TrimSpace(userID) != "" {
			next.ServeHTTP(w, r)
			return
		}
		if userID == "" || err != nil {
			profile, err := user.ProfileFromContext(r.Context())
			if err != nil || profile == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, err = uaa.CreateUser(profile.AccountName, origin, profile.AccountName, profile.Email)
			if err != nil || strings.TrimSpace(userID) == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			r = r.WithContext(session.ContextWithUserID(r.Context(), userID))
			session.UpdateSessionWithUserID(w, r, s, userID)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Authenticate guards access to protected resources by inspecting the user's token
func Authenticate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		token, err := session.TokenFromContext(req.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if token.Expiry.UTC().Sub(time.Now().UTC()) < 0 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		profile, err := user.ProfileFromContext(req.Context())
		if err != nil || strings.TrimSpace(profile.Email) == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// Authorize guards access to protected resources by inspecting the user's token
func Authorize(next http.Handler, domain string) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		profile, err := user.ProfileFromContext(req.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if strings.TrimSpace(domain) != "" && !strings.HasSuffix(strings.ToLower(profile.Email), domain) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// Secure ensures the current request is TLS secured, authenticated, and authorized
func Secure(next http.Handler, domain string, store sessions.Store) http.Handler {
	return ensureHTTPS(session.PopulateContext(Authenticate(Authorize(next, domain)), store))
}

// CallbackHandler handles Google redirection URI requests and adds the Google
// access token and Userinfoplus to the ctx. If authentication succeeds,
// handling delegates to the success handler, otherwise to the failure handler.
func CallbackHandler(config *oauth2.Config, fetcher user.Fetcher, success, failure http.Handler) http.Handler {
	wrappedSuccessHandler := func(config *oauth2.Config, f user.Fetcher, success, failure http.Handler) http.Handler {
		if failure == nil {
			failure = gologin.DefaultFailureHandler
		}
		fn := func(w http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			token, err := dgoauth2.TokenFromContext(ctx)
			if err != nil {
				ctx = gologin.WithError(ctx, err)
				failure.ServeHTTP(w, req.WithContext(ctx))
				return
			}

			profile, err := f.Profile(ctx, config, token)
			err = validateResponse(profile, err)
			if err != nil {
				ctx = gologin.WithError(ctx, err)
				failure.ServeHTTP(w, req.WithContext(ctx))
				return
			}
			ctx = user.WithProfile(ctx, profile)
			success.ServeHTTP(w, req.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}(config, fetcher, success, failure)
	return dgoauth2.CallbackHandler(config, wrappedSuccessHandler, failure)
}

// validateResponse returns an error if the given profile, raw
// http.Response, or error are unexpected. Returns nil if they are valid.
func validateResponse(profile *user.Profile, err error) error {
	if err != nil {
		return errors.Wrap(err, "unable to get Profile")
	}
	if profile == nil || profile.AccountName == "" {
		return errors.New("could not validate Profile")
	}
	return nil
}
