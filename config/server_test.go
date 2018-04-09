package config

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestServer(t *testing.T) {
	spec.Run(t, "Server", testServer, spec.Report(report.Terminal{}))
}

func testServer(t *testing.T, when spec.G, it spec.S) {
	reset := func() {
		os.Unsetenv("VCAP_APPLICATION")
		os.Unsetenv("VCAP_SERVICES")
		os.Unsetenv("PORT")

		os.Unsetenv("IGNITION_SCHEME")
		os.Unsetenv("IGNITION_DOMAIN")
		os.Unsetenv("IGNITION_PORT")
		os.Unsetenv("IGNITION_SERVE_PORT")
		os.Unsetenv("IGNITION_WEB_ROOT")
		os.Unsetenv("IGNITION_SESSION_SECRET")
	}
	it.Before(func() {
		RegisterTestingT(t)
		reset()
	})

	it.After(func() {
		reset()
	})

	when("not running on Cloud Foundry", func() {
		it.Before(func() {
			os.Unsetenv("VCAP_APPLICATION")
			os.Unsetenv("VCAP_SERVICES")
			os.Unsetenv("PORT")
		})

		when("the session secret is set", func() {
			it.Before(func() {
				os.Setenv("IGNITION_SESSION_SECRET", "test-session-secret")
			})

			it("errors if env vars have invalid contents", func() {
				os.Setenv("IGNITION_PORT", "%#%")
				s, err := NewServer()
				Expect(err).To(HaveOccurred())
				Expect(s).To(BeNil())
			})

			it("succeeds with default values when the environment is not set", func() {
				s, err := NewServer()
				Expect(err).NotTo(HaveOccurred())
				Expect(s).NotTo(BeNil())
				Expect(s.Scheme).To(Equal("http"))
				Expect(s.Domain).To(Equal("localhost"))
				Expect(s.Port).To(Equal(3000))
				Expect(s.ServePort).To(Equal(3000))
				Expect(s.WebRoot).To(ContainSubstring("dist"))
			})
		})

		when("the session secret is not set", func() {
			it("errors", func() {
				s, err := NewServer()
				Expect(err).To(HaveOccurred())
				Expect(s).To(BeNil())
			})
		})
	})

	when("running on Cloud Foundry", func() {
		it.Before(func() {
			os.Setenv("VCAP_APPLICATION", `{"cf_api": "https://api.run.pcfbeta.io","limits": {"fds": 16384},"application_name": "ignition","application_uris": ["ignition.pcfbeta.io"],"name": "ignition","space_name": "development","space_id": "test-space-id","uris": ["ignition.pcfbeta.io"],"users": null,"application_id": "test-app-id"}`)
			os.Setenv("VCAP_SERVICES", `{}`)
			os.Setenv("PORT", "54321")
		})

		when("the session secret is not set", func() {
			it("errors", func() {
				s, err := NewServer()
				Expect(err).To(HaveOccurred())
				Expect(s).To(BeNil())
			})
		})

		when("the session secret is set in the environment", func() {
			it.Before(func() {
				os.Setenv("IGNITION_SESSION_SECRET", "test-session-secret")
			})

			it("succeeds", func() {
				s, err := NewServer()
				Expect(err).NotTo(HaveOccurred())
				Expect(s).NotTo(BeNil())
				Expect(s.Scheme).To(Equal("https"))
				Expect(s.Domain).To(Equal("ignition.pcfbeta.io"))
				Expect(s.Port).To(Equal(443))
				Expect(s.ServePort).To(Equal(54321))
				Expect(s.WebRoot).NotTo(ContainSubstring("dist"))
				Expect(s.SessionSecret).To(Equal("test-session-secret"))
			})

			when("a domain is configured", func() {
				it.Before(func() {
					os.Setenv("IGNITION_DOMAIN", "ignition.example.com")
				})

				it.After(func() {
					os.Unsetenv("IGNITION_DOMAIN")
				})

				it("does not overwrite the IGNITION_DOMAIN value", func() {
					s, err := NewServer()
					Expect(err).NotTo(HaveOccurred())
					Expect(s).NotTo(BeNil())
					Expect(s.Scheme).To(Equal("https"))
					Expect(s.Domain).To(Equal("ignition.example.com"))
				})

				when("a domain is configured in ignition-creds", func() {
					it.Before(func() {
						os.Setenv("VCAP_SERVICES", `{
						  "user-provided": [
						    {
						      "name": "ignition-config",
						      "instance_name": "ignition-config",
						      "credentials": {
						        "domain": "ignition.example.net"
						      }
						    }
						  ]
						}`)
					})

					it("overwrites the IGNITION_DOMAIN value with the ignition-creds domain value", func() {
						s, err := NewServer()
						Expect(err).NotTo(HaveOccurred())
						Expect(s).NotTo(BeNil())
						Expect(s.Scheme).To(Equal("https"))
						Expect(s.Domain).To(Equal("ignition.example.net"))
					})
				})
			})

			when("there is no domain configured and no application URIs", func() {
				it.Before(func() {
					os.Setenv("VCAP_APPLICATION", `{"cf_api": "https://api.run.pcfbeta.io","limits": {"fds": 16384},"application_name": "ignition","application_uris": [],"name": "ignition","space_name": "development","space_id": "test-space-id","uris": [],"users": null,"application_id": "test-app-id"}`)
				})

				it("errors", func() {
					s, err := NewServer()
					Expect(err).To(HaveOccurred())
					Expect(s).To(BeNil())
				})
			})

			when("VCAP_APPLICATION contents are invalid", func() {
				it.Before(func() {
					os.Setenv("VCAP_APPLICATION", `*&^#%#`)
				})

				it("errors", func() {
					s, err := NewServer()
					Expect(err).To(HaveOccurred())
					Expect(s).To(BeNil())
				})
			})
		})

		when("the session secret is set in ignition-config", func() {
			it.Before(func() {
				os.Setenv("VCAP_SERVICES", `{"user-provided": [{
					"name": "ignition-config",
					"instance_name": "ignition-config",
					"credentials": {
						"session_secret": "test-config-session-secret"
					}}]}`)
			})

			it("succeeds", func() {
				s, err := NewServer()
				Expect(err).NotTo(HaveOccurred())
				Expect(s).NotTo(BeNil())
				Expect(s.SessionSecret).To(Equal("test-config-session-secret"))
			})

			when("the session secret is also set in the environment", func() {
				it.Before(func() {
					os.Setenv("IGNITION_SESSION_SECRET", "test-session-secret")
				})

				it("still uses the session secret from ignition-config", func() {
					s, err := NewServer()
					Expect(err).NotTo(HaveOccurred())
					Expect(s).NotTo(BeNil())
					Expect(s.SessionSecret).To(Equal("test-config-session-secret"))
				})
			})
		})
	})
}
