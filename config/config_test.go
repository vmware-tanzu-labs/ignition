package config

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/internal"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

type info struct {
	AuthorizationEndpoint  string `json:"authorization_endpoint"`
	TokenEndpoint          string `json:"token_endpoint"`
	AppSSHEndpoint         string `json:"app_ssh_endpoint"`
	DopplerLoggingEndpoint string `json:"doppler_logging_endpoint"`
	RoutingEndpoint        string `json:"routing_endpoint"`
}

func TestNew(t *testing.T) {
	spec.Run(t, "New", testNew, spec.Report(report.Terminal{}))
}

func testNew(t *testing.T, when spec.G, it spec.S) {
	reset := func() {
		os.Unsetenv("VCAP_APPLICATION")
		os.Unsetenv("VCAP_SERVICES")
		os.Unsetenv("PORT")

		// Server
		os.Unsetenv("IGNITION_SCHEME")
		os.Unsetenv("IGNITION_DOMAIN")
		os.Unsetenv("IGNITION_PORT")
		os.Unsetenv("IGNITION_SERVE_PORT")
		os.Unsetenv("IGNITION_WEB_ROOT")
		os.Unsetenv("IGNITION_SESSION_SECRET") // REQUIRED

		// Deployment
		os.Unsetenv("IGNITION_SYSTEM_DOMAIN")     // REQUIRED
		os.Unsetenv("IGNITION_UAA_ORIGIN")        // REQUIRED
		os.Unsetenv("IGNITION_API_CLIENT_ID")     // REQUIRED
		os.Unsetenv("IGNITION_API_CLIENT_SECRET") // REQUIRED

		// Experimenter
		os.Unsetenv("IGNITION_ORG_PREFIX")
		os.Unsetenv("IGNITION_QUOTA_NAME")
		os.Unsetenv("IGNITION_SPACE_NAME")

		// Authorizer
		os.Unsetenv("IGNITION_AUTH_VARIANT")
		os.Unsetenv("IGNITION_CLIENT_ID")         // REQUIRED
		os.Unsetenv("IGNITION_CLIENT_SECRET")     // REQUIRED
		os.Unsetenv("IGNITION_AUTH_URL")          // REQUIRED
		os.Unsetenv("IGNITION_AUTHORIZED_DOMAIN") // REQUIRED
	}

	var (
		s *httptest.Server
	)

	it.Before(func() {
		RegisterTestingT(t)
		reset()
	})

	it.After(func() {
		if s != nil {
			s.Close()
		}
	})

	when("not running on CF and all required env vars are set", func() {
		var (
			isoSegmentFile string
			quotaFile      string
			wellKnownFile  string
		)

		it.Before(func() {
			isoSegmentFile = "isolation-segments.json"
			quotaFile = "quota.json"
			wellKnownFile = "well-known.json"

			s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.Path, "well-known") {
					handler := internal.HandleTestdata(t, wellKnownFile, func() {})
					handler.ServeHTTP(w, r)
					return
				}
				if strings.Contains(r.URL.Path, "quota") {
					handler := internal.HandleTestdata(t, quotaFile, func() {})
					handler.ServeHTTP(w, r)
					return
				}
				if strings.Contains(r.URL.Path, "isolation_segments") {
					handler := internal.HandleTestdata(t, isoSegmentFile, func() {})
					handler.ServeHTTP(w, r)
					return
				}
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
						return
					})
					handler.ServeHTTP(w, r)
					return
				}
				if strings.Contains(r.URL.Path, "token") {
					handler := internal.HandleTestdata(t, "token.json", func() {})
					handler.ServeHTTP(w, r)
					return
				}
			}))
			os.Setenv("IGNITION_SESSION_SECRET", "test-ignition-session-secret")
			os.Setenv("IGNITION_SYSTEM_DOMAIN", s.URL)
			os.Setenv("IGNITION_UAA_ORIGIN", "test-ignition-uaa-origin")
			os.Setenv("IGNITION_API_CLIENT_ID", "test-ignition-api-client-id")
			os.Setenv("IGNITION_API_CLIENT_SECRET", "test-ignition-api-client-secret")
			os.Setenv("IGNITION_CLIENT_ID", "test-ignition-client-id")
			os.Setenv("IGNITION_CLIENT_SECRET", "test-ignition-client-secret")
			os.Setenv("IGNITION_AUTH_URL", s.URL)
			os.Setenv("IGNITION_AUTHORIZED_DOMAIN", "test-ignition-authorized-domain")
		})

		it("succeeds", func() {
			i, err := New()
			Expect(err).NotTo(HaveOccurred())
			Expect(i).NotTo(BeNil())
		})

		when("IGNITION_SESSION_SECRET is not set", func() {
			it.Before(func() {
				os.Unsetenv("IGNITION_SESSION_SECRET")
			})

			it("errors", func() {
				i, err := New()
				Expect(err).To(HaveOccurred())
				Expect(i).To(BeNil())
			})
		})

		when("IGNITION_AUTHORIZED_DOMAIN is not set", func() {
			it.Before(func() {
				os.Unsetenv("IGNITION_AUTHORIZED_DOMAIN")
			})

			it("errors", func() {
				i, err := New()
				Expect(err).To(HaveOccurred())
				Expect(i).To(BeNil())
			})
		})

		when("IGNITION_AUTH_URL is not set", func() {
			it.Before(func() {
				os.Unsetenv("IGNITION_AUTH_URL")
			})

			it("errors", func() {
				i, err := New()
				Expect(err).To(HaveOccurred())
				Expect(i).To(BeNil())
			})
		})

		when("IGNITION_CLIENT_ID is not set", func() {
			it.Before(func() {
				os.Unsetenv("IGNITION_CLIENT_ID")
			})

			it("errors", func() {
				i, err := New()
				Expect(err).To(HaveOccurred())
				Expect(i).To(BeNil())
			})
		})

		when("IGNITION_CLIENT_SECRET is not set", func() {
			it.Before(func() {
				os.Unsetenv("IGNITION_CLIENT_SECRET")
			})

			it("errors", func() {
				i, err := New()
				Expect(err).To(HaveOccurred())
				Expect(i).To(BeNil())
			})
		})

		when("IGNITION_SYSTEM_DOMAIN is not set", func() {
			it.Before(func() {
				os.Unsetenv("IGNITION_SYSTEM_DOMAIN")
			})

			it("errors", func() {
				i, err := New()
				Expect(err).To(HaveOccurred())
				Expect(i).To(BeNil())
			})
		})

		when("IGNITION_UAA_ORIGIN is not set", func() {
			it.Before(func() {
				os.Unsetenv("IGNITION_UAA_ORIGIN")
			})

			it("errors", func() {
				i, err := New()
				Expect(err).To(HaveOccurred())
				Expect(i).To(BeNil())
			})
		})

		when("IGNITION_API_CLIENT_ID is not set", func() {
			it.Before(func() {
				os.Unsetenv("IGNITION_API_CLIENT_ID")
			})

			it("errors", func() {
				i, err := New()
				Expect(err).To(HaveOccurred())
				Expect(i).To(BeNil())
			})
		})

		when("IGNITION_API_CLIENT_SECRET is not set", func() {
			it.Before(func() {
				os.Unsetenv("IGNITION_API_CLIENT_SECRET")
			})

			it("errors", func() {
				i, err := New()
				Expect(err).To(HaveOccurred())
				Expect(i).To(BeNil())
			})
		})

		when("the quota cannot be found", func() {
			it.Before(func() {
				quotaFile = "empty-quota.json"
			})

			it("errors", func() {
				i, err := New()
				Expect(err).To(HaveOccurred())
				Expect(i).To(BeNil())
			})
		})

		when("the well known metadata is invalid", func() {
			it.Before(func() {
				wellKnownFile = "invalid-well-known.json"
			})

			it("errors", func() {
				i, err := New()
				Expect(err).To(HaveOccurred())
				Expect(i).To(BeNil())
			})
		})
	})
}
