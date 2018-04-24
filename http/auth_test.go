package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dgoauth2 "github.com/dghubble/gologin/oauth2"
	"github.com/dghubble/sessions"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/http/session"
	"github.com/pivotalservices/ignition/http/session/sessionfakes"
	"github.com/pivotalservices/ignition/uaa/uaafakes"
	"github.com/pivotalservices/ignition/user"
	"github.com/pivotalservices/ignition/user/openid"
	"github.com/pivotalservices/ignition/user/openid/openidfakes"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"golang.org/x/oauth2"
)

type fakeSessionStore struct {
	session       *sessions.Session
	newCalled     bool
	getCalled     bool
	saveCalled    bool
	destroyCalled bool
}

func newFakeSessionStore() *fakeSessionStore {
	return &fakeSessionStore{}
}

func (f *fakeSessionStore) New(name string) *sessions.Session {
	config := &sessions.Config{
		Path:     "/",
		MaxAge:   3600 * 24 * 7,
		HTTPOnly: true,
	}
	session := sessions.NewSession(f, name)
	session.Config = config
	f.session = session
	f.newCalled = true
	return f.session
}

func (f *fakeSessionStore) Get(req *http.Request, name string) (*sessions.Session, error) {
	f.getCalled = true
	return f.session, nil
}

func (f *fakeSessionStore) Save(w http.ResponseWriter, session *sessions.Session) error {
	f.saveCalled = true
	f.session = session
	return nil
}

func (f *fakeSessionStore) Destroy(w http.ResponseWriter, name string) {
	f.destroyCalled = true
	f.session = nil
}

func TestAuthenticate(t *testing.T) {
	spec.Run(t, "Authenticate", testAuthenticate, spec.Report(report.Terminal{}))
}

func testAuthenticate(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	it("is unauthorized when there is no token", func() {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		Authenticate(nil).ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusUnauthorized))
	})

	it("is unauthorized when there is an expired token", func() {
		w := httptest.NewRecorder()
		t := &oauth2.Token{
			AccessToken:  "test-token",
			RefreshToken: "test-refresh-token",
			TokenType:    "bearer",
			Expiry:       time.Now().Add(-24 * time.Hour),
		}
		ctx := session.ContextWithToken(context.Background(), t)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		Authenticate(nil).ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusUnauthorized))
	})

	when("there is a valid token", func() {
		var (
			w   *httptest.ResponseRecorder
			req *http.Request
			t   *oauth2.Token
		)

		it.Before(func() {
			w = httptest.NewRecorder()
			t = &oauth2.Token{
				AccessToken:  "test-token",
				RefreshToken: "test-refresh-token",
				TokenType:    "bearer",
				Expiry:       time.Now().Add(24 * time.Hour),
			}
			ctx := session.ContextWithToken(context.Background(), t)
			req = httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		})

		when("there isn't a valid profile in the context", func() {
			it("is unauthorized", func() {
				Authenticate(nil).ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		when("there an empty profile is in the context", func() {
			it.Before(func() {
				profile := &user.Profile{}
				req = req.WithContext(user.WithProfile(req.Context(), profile))
			})

			it("is unauthorized", func() {
				Authenticate(nil).ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		when("there is a valid profile in the context", func() {
			it.Before(func() {
				profile := &user.Profile{
					Name:        "Test User",
					Email:       "test@example.net",
					AccountName: "corp\tester",
				}
				req = req.WithContext(user.WithProfile(req.Context(), profile))
			})

			it("calls the next handler", func() {
				called := false
				next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					called = true
					w.WriteHeader(http.StatusOK)
				})
				Authenticate(next).ServeHTTP(w, req)
				Expect(called).To(BeTrue())
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})
	})
}

func TestAuthorize(t *testing.T) {
	spec.Run(t, "Authorize", testAuthorize, spec.Report(report.Terminal{}))
}

func testAuthorize(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	when("there is a valid token", func() {
		var (
			w   *httptest.ResponseRecorder
			req *http.Request
			t   *oauth2.Token
		)

		it.Before(func() {
			w = httptest.NewRecorder()
			t = &oauth2.Token{
				AccessToken:  "test-token",
				RefreshToken: "test-refresh-token",
				TokenType:    "bearer",
				Expiry:       time.Now().Add(24 * time.Hour),
			}
			ctx := session.ContextWithToken(context.Background(), t)
			req = httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		})

		it("is unauthorized if there is no profile in the context", func() {
			Authorize(nil, "").ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})

		when("there is a valid profile in the context", func() {
			it.Before(func() {
				profile := &user.Profile{
					Name:        "Test User",
					Email:       "test@example.net",
					AccountName: "corp\tester",
				}
				req = req.WithContext(user.WithProfile(req.Context(), profile))
			})

			it("is forbidden if the user's email does not end with the domain", func() {
				Authorize(nil, "example.com").ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusForbidden))
			})

			it("calls the next handler if the user's email does end with the domain", func() {
				called := false
				next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					called = true
					w.WriteHeader(http.StatusOK)
				})
				Authorize(next, "example.net").ServeHTTP(w, req)
				Expect(called).To(BeTrue())
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})
	})
}

func TestCallbackHandler(t *testing.T) {
	spec.Run(t, "CallbackHandler", testCallbackHandler, spec.Report(report.Terminal{}))
}

func testCallbackHandler(t *testing.T, when spec.G, it spec.S) {
	var (
		s            *httptest.Server
		fakeVerifier *openidfakes.FakeVerifier
	)

	it.Before(func() {
		RegisterTestingT(t)
		s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(
				`{
					"access_token": "access0-token",
					"token_type": "bearer",
					"id_token": "id-token",
					"refresh_token": "refresh-token",
					"expires_in": 7199,
					"scope": "cloud_controller.admin cloud_controller.write doppler.firehose openid scim.read uaa.user cloud_controller.read password.write scim.write",
					"jti": "1234567890"
				}`))
		}))
		fakeVerifier = &openidfakes.FakeVerifier{}
		fakeVerifier.VerifyReturns(&openid.Claims{
			UserName: "testuser",
			Email:    "test@pivotal.io",
		}, nil)
	})

	it.After(func() {
		s.Close()
	})

	it("", func() {
		fetcher := &openid.Fetcher{
			Verifier: fakeVerifier,
		}
		succeeded := false
		failed := false
		var ctxActual context.Context
		success := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			succeeded = true
			ctxActual = r.Context()
		})
		failure := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			failed = true
			ctxActual = r.Context()
		})
		handler := CallbackHandler(&oauth2.Config{
			Endpoint: oauth2.Endpoint{
				AuthURL:  s.URL,
				TokenURL: s.URL,
			},
		}, fetcher, success, failure)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/oauth2?state=teststate&code=testcode", nil)
		ctx := context.Background()
		ctx = dgoauth2.WithState(ctx, "teststate")
		handler.ServeHTTP(w, r.WithContext(ctx))

		Expect(failed).To(BeFalse())
		Expect(succeeded).To(BeTrue())

		p, err := user.ProfileFromContext(ctxActual)
		Expect(err).NotTo(HaveOccurred())
		Expect(p.Email).To(Equal("test@pivotal.io"))
	})
}

func TestEnsureUser(t *testing.T) {
	spec.Run(t, "EnsureUser", testEnsureUser, spec.Report(report.Terminal{}))
}

func testEnsureUser(t *testing.T, when spec.G, it spec.S) {
	var (
		called           bool
		next             http.Handler
		handler          http.Handler
		w                *httptest.ResponseRecorder
		r                *http.Request
		uaa              *uaafakes.FakeAPI
		fakeSessionStore *sessionfakes.FakeStore
	)

	it.Before(func() {
		RegisterTestingT(t)

		called = false
		uaa = &uaafakes.FakeAPI{}
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})
		fakeSessionStore = &sessionfakes.FakeStore{}
		s := sessions.NewSession(fakeSessionStore, "ignition-test")
		fakeSessionStore.SaveReturns(nil)
		fakeSessionStore.GetReturns(s, nil)
		handler = ensureUser(next, uaa, "origin", fakeSessionStore)
		w = httptest.NewRecorder()
		r = httptest.NewRequest(http.MethodGet, "/", nil)
	})

	when("the user ID is present in the context", func() {
		it("calls the next handler", func() {
			ctx := session.ContextWithUserID(context.Background(), "test-user")
			handler.ServeHTTP(w, r.WithContext(ctx))
			Expect(called).To(BeTrue())
		})
	})

	when("the user ID not present in the context", func() {
		when("there is no profile in the context", func() {
			it("is unauthorized", func() {
				handler.ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		when("there is a profile in the context", func() {
			var ctx context.Context

			it.Before(func() {
				ctx = user.WithProfile(context.Background(), &user.Profile{
					Email:       "test@pivotal.io",
					AccountName: "testaccount",
					Name:        "test",
				})
			})

			it("creates the user and calls the next handler", func() {
				uaa.CreateUserReturns("test-user-id", nil)
				handler.ServeHTTP(w, r.WithContext(ctx))
				Expect(called).To(BeTrue())
				Expect(uaa.CreateUserCallCount()).To(Equal(1))
			})

			it("is unauthorized if the user cannot be created", func() {
				uaa.CreateUserReturns("", errors.New("test error"))
				handler.ServeHTTP(w, r.WithContext(ctx))
				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			})
		})
	})
}
