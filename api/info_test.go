package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bitly/go-simplejson"
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/api"
	"github.com/pivotalservices/ignition/cloudfoundry/cloudfoundryfakes"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestInfoHandler(t *testing.T) {
	spec.Run(t, "InfoHandler", testInfoHandler, spec.Report(report.Terminal{}))
}

func testInfoHandler(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	when("the ignition org count updates async in the background", func() {
		var handler http.Handler
		it.Before(func() {
			a := &cloudfoundryfakes.FakeAPI{}
			handler = api.InfoHandler(
				"Test Company", "Test Space", "ignition-quota-definition-guid", false, 100*time.Millisecond, a)

			// stub this out after the handler has initialized, the goroutine will update
			a.ListOrgsByQueryReturns([]cfclient.Org{
				cfclient.Org{
					Guid:                "4321",
					Name:                "orgprefix-joe",
					QuotaDefinitionGuid: "ignition-quota-definition-guid",
				},
				cfclient.Org{
					Guid:                "5432",
					Name:                "orgprefix-larry",
					QuotaDefinitionGuid: "ignition-quota-definition-guid",
				},
			}, nil)
		})

		it("returns the updated ignition org count", func() {
			Eventually(func() int {
				r := httptest.NewRecorder()
				handler.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/", nil))
				Expect(r.Code).To(Equal(http.StatusOK))
				j, err := simplejson.NewFromReader(r.Body)
				if err != nil {
					t.Errorf("Error while reading response JSON: %s", err)
				}
				return j.GetPath("IgnitionOrgCount").MustInt()
			}, "2s").Should(Equal(2))
		})
	})

	when("there are ignition and non-ignition orgs", func() {
		var handler http.Handler

		it.Before(func() {
			a := &cloudfoundryfakes.FakeAPI{}
			a.ListOrgsByQueryReturns([]cfclient.Org{
				cfclient.Org{
					Guid:                "1234",
					Name:                "some random org",
					QuotaDefinitionGuid: "other-quota-definition-guid",
				},
				cfclient.Org{
					Guid:                "4321",
					Name:                "orgprefix-joe",
					QuotaDefinitionGuid: "ignition-quota-definition-guid",
				},
				cfclient.Org{
					Guid:                "5432",
					Name:                "orgprefix-larry",
					QuotaDefinitionGuid: "ignition-quota-definition-guid",
				},
			}, nil)
			handler = api.InfoHandler("Test Company", "Test Space", "ignition-quota-definition-guid", false, time.Minute, a)
		})

		it("returns the configured company name, space name, and ignition org count", func() {
			r := httptest.NewRecorder()
			handler.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/", nil))
			Expect(r.Code).To(Equal(http.StatusOK))

			j, err := simplejson.NewFromReader(r.Body)
			if err != nil {
				t.Errorf("Error while reading response JSON: %s", err)
			}

			Expect(j.GetPath("CompanyName").MustString()).To(Equal("Test Company"))
			Expect(j.GetPath("ExperimentationSpaceName").MustString()).To(Equal("Test Space"))
			Expect(j.GetPath("IgnitionOrgCount").MustInt()).To(Equal(2))
		})
	})

	when("the cc api returns an error", func() {
		var handler http.Handler

		it.Before(func() {
			a := &cloudfoundryfakes.FakeAPI{}
			a.ListOrgsByQueryReturns([]cfclient.Org{}, errors.New("Some unknown CC API error"))
			handler = api.InfoHandler("Test Company", "Test Space", "orgprefix", false, time.Minute, a)
		})

		it("returns the configured company name, space name, and defaults the org count to 0", func() {
			r := httptest.NewRecorder()
			handler.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/", nil))
			Expect(r.Code).To(Equal(http.StatusOK))

			j, err := simplejson.NewFromReader(r.Body)
			if err != nil {
				t.Errorf("Error while reading response JSON: %s", err)
			}

			Expect(j.GetPath("CompanyName").MustString()).To(Equal("Test Company"))
			Expect(j.GetPath("ExperimentationSpaceName").MustString()).To(Equal("Test Space"))
			Expect(j.GetPath("IgnitionOrgCount").MustInt()).To(Equal(0))
		})
	})
}
