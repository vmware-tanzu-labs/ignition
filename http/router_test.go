package http

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/cloudfoundry/cloudfoundryfakes"
	"github.com/pivotalservices/ignition/config"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestAPI(t *testing.T) {
	spec.Run(t, "API", testAPI, spec.Report(report.Terminal{}))
}

func testAPI(t *testing.T, when spec.G, it spec.S) {
	var api *API
	it.Before(func() {
		RegisterTestingT(t)
		api = &API{
			Ignition: &config.Ignition{
				Authorizer: &config.Authorizer{},
				Deployment: &config.Deployment{
					CC: &cloudfoundryfakes.FakeAPI{},
				},
				Experimenter: &config.Experimenter{},
				Server:       &config.Server{},
			},
		}
	})

	it("returns the URI correctly", func() {
		api.Ignition.Server.Scheme = "https"
		api.Ignition.Server.Domain = "example.net"
		Expect(api.URI()).To(Equal("https://example.net"))
		api.Ignition.Server.Port = 1234
		Expect(api.URI()).To(Equal("https://example.net:1234"))
	})

	it("creates a valid router", func() {
		r := api.createRouter()
		Expect(r).NotTo(BeNil())
		assets := r.GetRoute("assets")
		Expect(assets).NotTo(BeNil())
		nonexistent := r.GetRoute("nonexistent")
		Expect(nonexistent).To(BeNil())
	})
}
