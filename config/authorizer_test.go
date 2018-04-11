package config

import (
	"bytes"
	"fmt"
	"io"
	"log"
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

func TestNewAuthorizer(t *testing.T) {
	spec.Run(t, "NewAuthorizer", testNewAuthorizer, spec.Report(report.Terminal{}))
}

func testNewAuthorizer(t *testing.T, when spec.G, it spec.S) {
	var (
		s             *httptest.Server
		called        bool
		responseCode  int
		wellKnownFile string
	)

	reset := func() {
		os.Unsetenv("VCAP_APPLICATION")
		os.Unsetenv("VCAP_SERVICES")
		os.Unsetenv("PORT")

		os.Unsetenv("IGNITION_AUTH_VARIANT")
		os.Unsetenv("IGNITION_AUTH_SCOPES")
		os.Unsetenv("IGNITION_CLIENT_ID")
		os.Unsetenv("IGNITION_CLIENT_SECRET")
		os.Unsetenv("IGNITION_AUTH_URL")
		os.Unsetenv("IGNITION_AUTHORIZED_DOMAIN")
	}

	it.Before(func() {
		RegisterTestingT(t)
		reset()
		responseCode = http.StatusOK
		called = false
		wellKnownFile = "well-known.json"
		s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			wellKnown := internal.StringFromTestdata(t, wellKnownFile)
			wellKnown = strings.Replace(wellKnown, "{{url}}", s.URL, -1)
			w.WriteHeader(responseCode)
			io.Copy(w, bytes.NewReader([]byte(wellKnown)))
		}))

	})

	it.After(func() {
		reset()
		s.Close()
	})

	when("there are missing required variables", func() {
		it("errors", func() {
			a, err := NewAuthorizer("ignition-config")
			Expect(err).To(HaveOccurred())
			Expect(a).To(BeNil())
		})
	})

	when("all the required variables are set", func() {
		it.Before(func() {
			os.Setenv("IGNITION_AUTH_VARIANT", "test-ignition-auth-variant")
			os.Setenv("IGNITION_CLIENT_ID", "test-ignition-client-id")
			os.Setenv("IGNITION_CLIENT_SECRET", "test-ignition-client-secret")
			os.Setenv("IGNITION_AUTH_URL", s.URL)
			os.Setenv("IGNITION_AUTHORIZED_DOMAIN", "test-ignition-authorized-domain")
		})

		it("succeeds", func() {
			a, err := NewAuthorizer("ignition-config")
			Expect(err).NotTo(HaveOccurred())
			Expect(a).NotTo(BeNil())
			Expect(a.Variant).To(Equal("test-ignition-auth-variant"))
			Expect(a.ClientID).To(Equal("test-ignition-client-id"))
			Expect(a.ClientSecret).To(Equal("test-ignition-client-secret"))
			Expect(a.URL).To(Equal(s.URL))
			Expect(a.Domain).To(Equal("test-ignition-authorized-domain"))
		})

		when("there is an error fetching the well known metadata", func() {
			it.Before(func() {
				os.Setenv("IGNITION_AUTH_URL", "test^^#://$%&@")
			})

			it("errors", func() {
				a, err := NewAuthorizer("ignition-config")
				Expect(err).To(HaveOccurred())
				Expect(a).To(BeNil())
			})
		})

		when("the well known metadata is invalid", func() {
			it.Before(func() {
				wellKnownFile = "invalid-well-known.json"
			})

			it("errors", func() {
				a, err := NewAuthorizer("ignition-config")
				Expect(err).To(HaveOccurred())
				Expect(a).To(BeNil())
				log.Println(err)
			})
		})

		when("the well known metadata status is not http.StatusOK", func() {
			it.Before(func() {
				responseCode = http.StatusAccepted
			})

			it("errors", func() {
				a, err := NewAuthorizer("ignition-config")
				Expect(err).To(HaveOccurred())
				Expect(a).To(BeNil())
			})
		})

		it("errors when the client id is empty", func() {
			os.Unsetenv("IGNITION_CLIENT_ID")
			a, err := NewAuthorizer("ignition-config")
			Expect(err).To(HaveOccurred())
			Expect(a).To(BeNil())
		})

		it("errors when the client secret is empty", func() {
			os.Unsetenv("IGNITION_CLIENT_SECRET")
			a, err := NewAuthorizer("ignition-config")
			Expect(err).To(HaveOccurred())
			Expect(a).To(BeNil())
		})

		it("errors when the auth url is empty", func() {
			os.Unsetenv("IGNITION_AUTH_URL")
			a, err := NewAuthorizer("ignition-config")
			Expect(err).To(HaveOccurred())
			Expect(a).To(BeNil())
		})

		it("errors when the authorized domain is empty", func() {
			os.Unsetenv("IGNITION_AUTHORIZED_DOMAIN")
			a, err := NewAuthorizer("ignition-config")
			Expect(err).To(HaveOccurred())
			Expect(a).To(BeNil())
		})
	})

	when("running in CF", func() {
		it.Before(func() {
			os.Setenv("VCAP_APPLICATION", `{"cf_api": "https://api.run.pcfbeta.io","limits": {"fds": 16384},"application_name": "ignition","application_uris": ["ignition.pcfbeta.io"],"name": "ignition","space_name": "development","space_id": "test-space-id","uris": ["ignition.pcfbeta.io"],"users": null,"application_id": "test-app-id"}`)
			os.Setenv("VCAP_SERVICES", `{}`)
		})

		when("the VCAP_APPLICATION is invalid", func() {
			it.Before(func() {
				os.Setenv("VCAP_APPLICATION", "^&%#&^%")
			})

			it("errors", func() {
				a, err := NewAuthorizer("ignition-config")
				Expect(err).To(HaveOccurred())
				Expect(a).To(BeNil())
			})
		})

		when("there is no ignition-config service instance", func() {
			it("errors", func() {
				a, err := NewAuthorizer("ignition-config")
				Expect(err).To(HaveOccurred())
				Expect(a).To(BeNil())
			})
		})

		when("the ignition-config service instance exists", func() {
			it.Before(func() {
				os.Setenv("VCAP_SERVICES", fmt.Sprintf(`{
					"user-provided": [
						{
							"name": "ignition-config",
							"instance_name": "ignition-config",
							"credentials": {
								"auth_variant": "test-service-auth-variant",
								"auth_scopes": "testscope,anotherscope",
								"authorized_domain": "test-service-authorized-domain",
								"auth_url": "%s",
								"client_id": "test-service-client-id",
								"client_secret": "test-service-client-secret",
								"auth_servicename": "test-service-client-servicename"
							}
						}
					]
				}`, s.URL))
			})

			it("succeeds", func() {
				a, err := NewAuthorizer("ignition-config")
				Expect(err).NotTo(HaveOccurred())
				Expect(a).NotTo(BeNil())
				Expect(a.Variant).To(Equal("test-service-auth-variant"))
				Expect(a.Scopes[0]).To(Equal("testscope"))
				Expect(a.Scopes[1]).To(Equal("anotherscope"))
				Expect(a.ClientID).To(Equal("test-service-client-id"))
				Expect(a.ClientSecret).To(Equal("test-service-client-secret"))
				Expect(a.URL).To(Equal(s.URL))
				Expect(a.Domain).To(Equal("test-service-authorized-domain"))
				Expect(a.ServiceName).To(Equal("test-service-client-servicename"))
			})

			when("the variant is p-identity and there is no identity service", func() {
				it.Before(func() {
					os.Setenv("VCAP_SERVICES", `{
						"user-provided": [
							{
								"name": "ignition-config",
								"instance_name": "ignition-config",
								"credentials": {
									"auth_variant": "p-identity",
									"authorized_domain": "test-service-authorized-domain",
									"auth_url": "test-service-auth-url",
									"client_id": "test-service-client-id",
									"client_secret": "test-service-client-secret"
								}
							}
						]
					}`)
				})

				when("there is no identity service", func() {
					it("errors", func() {
						a, err := NewAuthorizer("ignition-config")
						Expect(err).To(HaveOccurred())
						Expect(a).To(BeNil())
					})
				})
			})

			when("the variant is p-identity and there is an identity service", func() {
				var vcapServices string

				it.Before(func() {
					vcapServices = fmt.Sprintf(`{
						"p-identity": [
							{
								"name": "ignition-identity",
								"instance_name": "ignition-identity",
								"binding_name": null,
								"credentials": {
									"auth_domain": "%s",
									"client_secret": "test-identity-client-secret",
									"client_id": "test-identity-client-id"
								},
								"syslog_drain_url": null,
								"volume_mounts": [],
								"label": "p-identity",
								"provider": null,
								"plan": "ignition",
								"tags": []
							}
						],
						"user-provided": [
							{
								"name": "ignition-config",
								"instance_name": "ignition-config",
								"credentials": {
									"auth_variant": "p-identity",
									"authorized_domain": "test-service-authorized-domain",
									"auth_url": "test-service-auth-url",
									"client_id": "test-service-client-id",
									"client_secret": "test-service-client-secret"
								}
							}
						]
					}`, s.URL)
					os.Setenv("VCAP_SERVICES", vcapServices)
				})

				it.After(func() {
					s.Close()
				})

				it("succeeds", func() {
					a, err := NewAuthorizer("ignition-config")
					Expect(err).NotTo(HaveOccurred())
					Expect(a).NotTo(BeNil())
				})

				when("the service is missing the auth_domain", func() {
					it.Before(func() {
						d := fmt.Sprintf("\"auth_domain\": \"%s\",", s.URL)
						os.Setenv("VCAP_SERVICES", strings.Replace(vcapServices, d, "", -1))
					})

					it("errors", func() {
						a, err := NewAuthorizer("ignition-config")
						Expect(err).To(HaveOccurred())
						Expect(a).To(BeNil())
					})
				})

				when("the service is missing the client_id", func() {
					it.Before(func() {
						d := strings.Replace(vcapServices, "\"client_id\": \"test-identity-client-id\"", "", -1)
						d = strings.Replace(d, "-client-secret\",", "-client-secret\"", -1)
						os.Setenv("VCAP_SERVICES", d)
					})

					it("errors", func() {
						a, err := NewAuthorizer("ignition-config")
						Expect(err).To(HaveOccurred())
						Expect(a).To(BeNil())
					})
				})

				when("the service is missing the client_secret", func() {
					it.Before(func() {
						os.Setenv("VCAP_SERVICES", strings.Replace(vcapServices, "\"client_secret\": \"test-identity-client-secret\",", "", -1))
					})

					it("errors", func() {
						a, err := NewAuthorizer("ignition-config")
						Expect(err).To(HaveOccurred())
						Expect(a).To(BeNil())
					})
				})
			})
		})
	})
}
