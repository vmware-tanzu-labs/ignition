package config

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDeployment(t *testing.T) {
	spec.Run(t, "Deployment", testDeployment, spec.Report(report.Terminal{}))
}

func testDeployment(t *testing.T, when spec.G, it spec.S) {
	reset := func() {
		os.Unsetenv("VCAP_APPLICATION")
		os.Unsetenv("VCAP_SERVICES")
		os.Unsetenv("PORT")

		os.Unsetenv("IGNITION_SYSTEM_DOMAIN")
		os.Unsetenv("IGNITION_UAA_ORIGIN")
		os.Unsetenv("IGNITION_API_CLIENT_ID")
		os.Unsetenv("IGNITION_API_CLIENT_SECRET")
		os.Unsetenv("IGNITION_SKIP_TLS_VALIDATION")
	}

	it.Before(func() {
		RegisterTestingT(t)
		reset()
	})

	it.After(func() {
		reset()
	})

	when("not running on Cloud Foundry", func() {
		it("errors if required env vars arent set", func() {
			d, err := NewDeployment("ignition-config")
			Expect(err).To(HaveOccurred())
			Expect(d).To(BeNil())
		})

		it("returns the correct URL", func() {
			d := &Deployment{
				SystemDomain: "run.example.com",
			}
			Expect(d.URL("").String()).To(Equal("https://run.example.com"))
			Expect(d.URL("apps").String()).To(Equal("https://apps.run.example.com"))
		})

		it("allows http to be used", func() {
			d := &Deployment{
				SystemDomain: "http://run.example.com",
			}
			Expect(d.URL("").String()).To(Equal("http://run.example.com"))
			Expect(d.URL("apps").String()).To(Equal("http://apps.run.example.com"))
		})

		it("is nil if the system domain is bad", func() {
			d := &Deployment{
				SystemDomain: "(*#&^@%$&%)",
			}
			Expect(d.URL("")).To(BeNil())
		})

		it("parses the system domain correctly", func() {
			d := &Deployment{
				SystemDomain: "run.example.com",
			}
			d.ParseSystemDomain()
			Expect(d.SystemDomain).To(Equal("https://run.example.com"))
			Expect(d.AppsURL).To(Equal("https://apps.run.example.com"))
			Expect(d.APIURL).To(Equal("https://api.run.example.com"))
			Expect(d.UAAURL).To(Equal("https://login.run.example.com"))
		})

		when("all required environment variables are set", func() {
			it.Before(func() {
				os.Setenv("IGNITION_SYSTEM_DOMAIN", "run.example.com")
				os.Setenv("IGNITION_UAA_ORIGIN", "okta")
				os.Setenv("IGNITION_API_CLIENT_ID", "test-client-id")
				os.Setenv("IGNITION_API_CLIENT_SECRET", "test-client-secret")
			})

			it("succeeds", func() {
				d, err := NewDeployment("ignition-config")
				Expect(err).NotTo(HaveOccurred())
				Expect(d).NotTo(BeNil())
				Expect(d.SystemDomain).To(Equal("https://run.example.com"))
				Expect(d.UAAOrigin).To(Equal("okta"))
				Expect(d.ClientID).To(Equal("test-client-id"))
				Expect(d.ClientSecret).To(Equal("test-client-secret"))
			})

			it("can generate an oauth2.Config", func() {
				d, err := NewDeployment("ignition-config")
				Expect(err).NotTo(HaveOccurred())
				Expect(d).NotTo(BeNil())
				c := d.Config()
				Expect(c).NotTo(BeNil())
				Expect(c.ClientID).To(Equal("test-client-id"))
				Expect(c.ClientSecret).To(Equal("test-client-secret"))
				Expect(c.TokenURL).To(Equal("https://login.run.example.com/oauth/token"))
				Expect(c.Scopes).To(HaveLen(3))
				Expect(c.Scopes).To(ConsistOf("cloud_controller.admin", "scim.write", "scim.read"))
			})

			when("skip tls validation is true", func() {
				it.Before(func() {
					os.Setenv("IGNITION_SKIP_TLS_VALIDATION", "true")
				})

				it("configures the cf client to skip ssl validation", func() {
					d, err := NewDeployment("ignition-config")
					Expect(err).ToNot(HaveOccurred())
					Expect(d.SkipTLSValidation).To(BeTrue())
				})
			})

			when("the system domain is empty", func() {
				it.Before(func() {
					os.Setenv("IGNITION_SYSTEM_DOMAIN", "")
				})

				it("errors", func() {
					d, err := NewDeployment("ignition-config")
					Expect(err).To(HaveOccurred())
					Expect(d).To(BeNil())
				})
			})

			when("the uaa origin is empty", func() {
				it.Before(func() {
					os.Setenv("IGNITION_UAA_ORIGIN", "")
				})

				it("errors", func() {
					d, err := NewDeployment("ignition-config")
					Expect(err).To(HaveOccurred())
					Expect(d).To(BeNil())
				})
			})

			when("the api client id is empty", func() {
				it.Before(func() {
					os.Setenv("IGNITION_API_CLIENT_ID", "")
				})

				it("errors", func() {
					d, err := NewDeployment("ignition-config")
					Expect(err).To(HaveOccurred())
					Expect(d).To(BeNil())
				})
			})

			when("the api client secret is empty", func() {
				it.Before(func() {
					os.Setenv("IGNITION_API_CLIENT_SECRET", "")
				})

				it("errors", func() {
					d, err := NewDeployment("ignition-config")
					Expect(err).To(HaveOccurred())
					Expect(d).To(BeNil())
				})
			})
		})
	})

	when("running on Cloud Foundry", func() {
		it.Before(func() {
			os.Setenv("VCAP_APPLICATION", `{"cf_api": "https://api.run.pcfbeta.io","limits": {"fds": 16384},"application_name": "ignition","application_uris": ["ignition.pcfbeta.io"],"name": "ignition","space_name": "development","space_id": "test-space-id","uris": ["ignition.pcfbeta.io"],"users": null,"application_id": "test-app-id"}`)
			os.Setenv("VCAP_SERVICES", `{}`)
		})

		it("fails when VCAP_APPLICATION is invalid", func() {
			os.Setenv("VCAP_APPLICATION", "^&%#&^%")
			d, err := NewDeployment("ignition-config")
			Expect(err).To(HaveOccurred())
			Expect(d).To(BeNil())
		})

		it("fails when a non-existent service is requested", func() {
			d, err := NewDeployment("ignition-config")
			Expect(err).To(HaveOccurred())
			Expect(d).To(BeNil())
		})

		when("all the fields are configured in the service", func() {
			it.Before(func() {
				os.Setenv("VCAP_SERVICES", `{
				  "user-provided": [
				    {
				      "name": "ignition-config",
				      "instance_name": "ignition-config",
				      "binding_name": null,
				      "credentials": {
				        "system_domain": "run.example.com",
				        "uaa_origin": "okta",
				        "api_client_id": "test-client-id",
				        "api_client_secret": "test-client-secret"
				      },
				      "syslog_drain_url": "",
				      "volume_mounts": [],
				      "label": "user-provided",
				      "tags": []
				    }
				  ]
				}`)
			})

			it("succeeds", func() {
				d, err := NewDeployment("ignition-config")
				Expect(err).NotTo(HaveOccurred())
				Expect(d).NotTo(BeNil())
				Expect(d.SystemDomain).To(Equal("https://run.example.com"))
				Expect(d.UAAOrigin).To(Equal("okta"))
				Expect(d.ClientID).To(Equal("test-client-id"))
				Expect(d.ClientSecret).To(Equal("test-client-secret"))
			})
		})

		when("skip tls validation is true", func() {
			it.Before(func() {
				os.Setenv("VCAP_SERVICES", `{
				  "user-provided": [
				    {
				      "name": "ignition-config",
				      "instance_name": "ignition-config",
				      "binding_name": null,
				      "credentials": {
				        "system_domain": "run.example.com",
				        "uaa_origin": "okta",
				        "api_client_id": "test-client-id",
				        "api_client_secret": "test-client-secret",
								"skip_tls_validation": "true"
				      },
				      "syslog_drain_url": "",
				      "volume_mounts": [],
				      "label": "user-provided",
				      "tags": []
				    }
				  ]
				}`)
			})

			it("configures the cf client to skip ssl validation", func() {
				d, err := NewDeployment("ignition-config")
				Expect(err).ToNot(HaveOccurred())
				Expect(d.SkipTLSValidation).To(BeTrue())
			})
		})

		when("the service is missing the system domain", func() {
			it.Before(func() {
				os.Setenv("VCAP_SERVICES", `{
				  "user-provided": [
				    {
				      "name": "ignition-config",
				      "instance_name": "ignition-config",
				      "binding_name": null,
				      "credentials": {
				        "uaa_origin": "okta",
				        "api_client_id": "test-client-id",
				        "api_client_secret": "test-client-secret"
				      },
				      "syslog_drain_url": "",
				      "volume_mounts": [],
				      "label": "user-provided",
				      "tags": []
				    }
				  ]
				}`)
			})

			it("errors", func() {
				d, err := NewDeployment("ignition-config")
				Expect(err).To(HaveOccurred())
				Expect(d).To(BeNil())
			})
		})

		when("the service is missing the uaa origin", func() {
			it.Before(func() {
				os.Setenv("VCAP_SERVICES", `{
				  "user-provided": [
				    {
				      "name": "ignition-config",
				      "instance_name": "ignition-config",
				      "binding_name": null,
				      "credentials": {
								"system_domain": "run.example.com",
				        "api_client_id": "test-client-id",
				        "api_client_secret": "test-client-secret"
				      },
				      "syslog_drain_url": "",
				      "volume_mounts": [],
				      "label": "user-provided",
				      "tags": []
				    }
				  ]
				}`)
			})

			it("errors", func() {
				d, err := NewDeployment("ignition-config")
				Expect(err).To(HaveOccurred())
				Expect(d).To(BeNil())
			})
		})

		when("the service is missing the api client id", func() {
			it.Before(func() {
				os.Setenv("VCAP_SERVICES", `{
				  "user-provided": [
				    {
				      "name": "ignition-config",
				      "instance_name": "ignition-config",
				      "binding_name": null,
				      "credentials": {
								"system_domain": "run.example.com",
								"uaa_origin": "okta",
				        "api_client_secret": "test-client-secret"
				      },
				      "syslog_drain_url": "",
				      "volume_mounts": [],
				      "label": "user-provided",
				      "tags": []
				    }
				  ]
				}`)
			})

			it("errors", func() {
				d, err := NewDeployment("ignition-config")
				Expect(err).To(HaveOccurred())
				Expect(d).To(BeNil())
			})
		})

		when("the service is missing the api client secret", func() {
			it.Before(func() {
				os.Setenv("VCAP_SERVICES", `{
				  "user-provided": [
				    {
				      "name": "ignition-config",
				      "instance_name": "ignition-config",
				      "binding_name": null,
				      "credentials": {
								"system_domain": "run.example.com",
								"uaa_origin": "okta",
								"api_client_id": "test-client-id"
				      },
				      "syslog_drain_url": "",
				      "volume_mounts": [],
				      "label": "user-provided",
				      "tags": []
				    }
				  ]
				}`)
			})

			it("errors", func() {
				d, err := NewDeployment("ignition-config")
				Expect(err).To(HaveOccurred())
				Expect(d).To(BeNil())
			})
		})
	})
}
