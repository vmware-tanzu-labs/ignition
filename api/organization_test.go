package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitly/go-simplejson"
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/api"
	"github.com/pivotalservices/ignition/cloudfoundry/cloudfoundryfakes"
	"github.com/pivotalservices/ignition/http/session"
	"github.com/pivotalservices/ignition/user"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestHandler(t *testing.T) {
	spec.Run(t, "Handler", testHandler, spec.Report(report.Terminal{}))
}

func testHandler(t *testing.T, when spec.G, it spec.S) {
	var r *http.Request
	var w *httptest.ResponseRecorder
	var c *cloudfoundryfakes.FakeAPI

	it.Before(func() {
		RegisterTestingT(t)
		w = httptest.NewRecorder()
		r = httptest.NewRequest(http.MethodGet, "/", nil)
		c = &cloudfoundryfakes.FakeAPI{}
	})

	when("there is no profile in the context", func() {
		it("is not found", func() {
			r = httptest.NewRequest(http.MethodGet, "/", nil)
			api.OrganizationHandler("http://example.net", "ignition", "test-quota-id", "test-iso-segment-id", "playground", c).ServeHTTP(w, r)
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	when("there is a profile in the context but no user id", func() {
		it("is not found", func() {
			r = httptest.NewRequest(http.MethodGet, "/", nil)
			profile := &user.Profile{
				AccountName: "testuser@test.com",
			}
			r = r.WithContext(user.WithProfile(r.Context(), profile))
			api.OrganizationHandler("http://example.net", "ignition", "test-quota-id", "test-iso-segment-id", "playground", c).ServeHTTP(w, r)
			Expect(w.Code).To(Equal(http.StatusNotFound))
		})
	})

	when("there is a profile and a user id in the context", func() {
		it.Before(func() {
			r = httptest.NewRequest(http.MethodGet, "/", nil)
			profile := &user.Profile{
				AccountName: "testuser@test.com",
			}
			r = r.WithContext(user.WithProfile(session.ContextWithUserID(r.Context(), "test-user-id"), profile))
		})

		when("orgs cannot be retrieved", func() {
			it.Before(func() {
				c.ListOrgsByQueryReturns(nil, errors.New("test error"))
			})

			it("is not found", func() {
				api.OrganizationHandler("http://example.net", "ignition", "test-quota-id", "test-iso-segment-id", "playground", c).ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		when("there are no orgs for the user", func() {
			it.Before(func() {
				c.ListOrgsByQueryReturns(nil, nil)
			})

			it("creates the org", func() {
				c.CreateOrgReturns(cfclient.Org{
					Guid:                        "test-org-guid",
					Name:                        "ignition-testuser",
					QuotaDefinitionGuid:         "test-quota-id",
					DefaultIsolationSegmentGuid: "test-iso-segment-id",
				}, nil)
				api.OrganizationHandler("http://example.net", "ignition", "test-quota-id", "test-iso-segment-id", "playground", c).ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(ContainSubstring("ignition-testuser"))
			})
		})

		when("there are multiple orgs for the user", func() {
			it.Before(func() {
				c.ListOrgsByQueryReturns([]cfclient.Org{
					cfclient.Org{
						Guid:                        "test-org-2",
						Name:                        "ignition-testuser1",
						QuotaDefinitionGuid:         "ignition-quota2-id",
						DefaultIsolationSegmentGuid: "test-iso-segment-id",
						CreatedAt:                   "created-at",
						UpdatedAt:                   "updated-at",
					},
					cfclient.Org{
						Guid:                        "test-org-1",
						Name:                        "ignition-testuser",
						QuotaDefinitionGuid:         "ignition-quota-id",
						DefaultIsolationSegmentGuid: "test-iso-segment-id",
						CreatedAt:                   "created-at",
						UpdatedAt:                   "updated-at",
					},
				}, nil)
			})

			it("selects the correct org when there is a name match", func() {
				api.OrganizationHandler("http://example.net", "ignition", "test-quota2-id", "test-iso-segment-id", "playground", c).ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(ContainSubstring("test-org-1"))
			})

			when("creating an org succeeds", func() {
				it.Before(func() {
					c.CreateOrgReturns(cfclient.Org{
						Guid:                        "test-org-guid",
						Name:                        "ignition1-testuser",
						QuotaDefinitionGuid:         "test-quota2-id",
						DefaultIsolationSegmentGuid: "test-iso-segment-id",
					}, nil)
				})

				it("creates the org when there is no name or quota match", func() {
					api.OrganizationHandler("http://example.net", "ignition1", "test-quota2-id", "test-iso-segment-id", "playground", c).ServeHTTP(w, r)
					Expect(w.Code).To(Equal(http.StatusOK))
					j, err := simplejson.NewFromReader(w.Body)
					if err != nil {
						t.Errorf("Error while reading response JSON: %s", err)
					}
					Expect(j.GetPath("guid").MustString()).To(Equal("test-org-guid"))
					Expect(j.GetPath("name").MustString()).To(Equal("ignition1-testuser"))
					Expect(j.GetPath("quota_definition_guid").MustString()).To(Equal("test-quota2-id"))
					Expect(j.GetPath("default_isolation_segment_guid").MustString()).To(Equal("test-iso-segment-id"))
				})
			})

			when("creating an org fails", func() {
				it.Before(func() {
					c.CreateOrgReturns(cfclient.Org{}, errors.New("test error"))
				})

				it("is not found", func() {
					api.OrganizationHandler("http://example.net", "ignition1", "test-quota2-id", "test-iso-segment-id", "playground", c).ServeHTTP(w, r)
					Expect(w.Code).To(Equal(http.StatusNotFound))
				})
			})

			it("selects the correct org when there is a quota match (but not a name match)", func() {
				api.OrganizationHandler("http://example.net", "ignition2", "ignition-quota-id", "test-iso-segment-id", "playground", c).ServeHTTP(w, r)
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Body.String()).To(ContainSubstring("test-org-1"))
			})
		})
	})
}

func TestOrgName(t *testing.T) {
	spec.Run(t, "OrgName", testOrgName, spec.Report(report.Terminal{}))
}

func testOrgName(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	it("works with email addresses", func() {
		Expect(api.OrganizationName("ignition", "test@example.net")).To(Equal("ignition-test"))
		Expect(api.OrganizationName("igNiTion", "tEsT@example.net")).To(Equal("ignition-test"))
	})

	it("works with domain accounts", func() {
		Expect(api.OrganizationName("ignition", "corp\\test")).To(Equal("ignition-test"))
		Expect(api.OrganizationName("igNiTion", "corp\\tEsT")).To(Equal("ignition-test"))
	})

	it("works with plain accounts", func() {
		Expect(api.OrganizationName("ignition", "test")).To(Equal("ignition-test"))
		Expect(api.OrganizationName("igNiTion", "tEsT")).To(Equal("ignition-test"))
	})
}
