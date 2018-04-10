package cloudfoundry

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/internal"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"golang.org/x/oauth2"
)

type fakeTokenSource struct{}

func (t *fakeTokenSource) Token() (*oauth2.Token, error) {
	return nil, errors.New("test error")
}

type fakeValidTokenSource struct{}

func (t *fakeValidTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken:  "test-token",
		TokenType:    "bearer",
		Expiry:       time.Now().Add(24 * time.Hour),
		RefreshToken: "test-refresh-token",
	}, nil
}

type info struct {
	AuthorizationEndpoint  string `json:"authorization_endpoint"`
	TokenEndpoint          string `json:"token_endpoint"`
	AppSSHEndpoint         string `json:"app_ssh_endpoint"`
	DopplerLoggingEndpoint string `json:"doppler_logging_endpoint"`
	RoutingEndpoint        string `json:"routing_endpoint"`
}

func TestClient(t *testing.T) {
	spec.Run(t, "Client", testClient, spec.Report(report.Terminal{}))
}

func testClient(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	when("authentication is required", func() {
		var (
			s         *httptest.Server
			c         *Client
			tokenFile string
		)

		it.Before(func() {
			tokenFile = "token.json"
			s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.Path, "info") {
					handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						i := info{
							AuthorizationEndpoint:  s.URL,
							TokenEndpoint:          s.URL,
							DopplerLoggingEndpoint: s.URL,
							RoutingEndpoint:        s.URL,
							AppSSHEndpoint:         s.URL,
						}
						json.NewEncoder(w).Encode(&i)
					})
					handler.ServeHTTP(w, r)
					return
				}
				if strings.Contains(r.URL.Path, "token") {
					handler := internal.HandleTestdata(t, tokenFile, func() {})
					handler.ServeHTTP(w, r)
					return
				}
			}))

			c = &Client{
				Config: &cfclient.Config{
					ApiAddress: s.URL,
				},
			}
		})

		it.After(func() {
			s.Close()
		})

		it("can authenticate", func() {
			err := c.authenticate()
			Expect(err).To(BeNil())
		})

		it("sets cf", func() {
			Expect(c.CF).To(BeNil())
			err := c.checkAuthentication()
			Expect(err).To(BeNil())
			Expect(c.CF).NotTo(BeNil())
		})

		when("cf is set but the token source errors", func() {
			it.Before(func() {
				c.CF = &cfclient.Client{
					Config: cfclient.Config{
						TokenSource: &fakeTokenSource{},
					},
				}
			})

			it("authenticates", func() {
				err := c.checkAuthentication()
				Expect(err).To(BeNil())
			})
		})

		when("cf is set and the token source returns a token", func() {
			it.Before(func() {
				c.CF = &cfclient.Client{
					Config: cfclient.Config{
						TokenSource: &fakeValidTokenSource{},
					},
				}
			})

			it("authenticates", func() {
				err := c.checkAuthentication()
				Expect(err).To(BeNil())
			})
		})

		when("the token is bad", func() {
			it.Before(func() {
				tokenFile = "invalid-token.json"
			})

			it("errors", func() {
				err := c.authenticate()
				Expect(err).NotTo(BeNil())
			})
		})
	})
}
