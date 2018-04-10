package cloudfoundry_test

import (
	"errors"
	"testing"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/cloudfoundry"
	"github.com/pivotalservices/ignition/cloudfoundry/cloudfoundryfakes"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var badCreator = func(config *cfclient.Config) (client cloudfoundry.API, err error) {
	return nil, errors.New("test-error")
}

func TestClientExternal(t *testing.T) {
	spec.Run(t, "ClientExternal", testClientExternal, spec.Report(report.Terminal{}))
}

func testClientExternal(t *testing.T, when spec.G, it spec.S) {
	var (
		c  *cloudfoundry.Client
		cf *cloudfoundryfakes.FakeAPI
	)

	it.Before(func() {
		RegisterTestingT(t)
		cf = &cloudfoundryfakes.FakeAPI{}
		c = &cloudfoundry.Client{
			CF: cf,
			Creator: func(config *cfclient.Config) (client cloudfoundry.API, err error) {
				return cf, nil
			},
		}
	})

	it("passes through to GetToken()", func() {
		cf.GetTokenReturns("test-token", nil)
		t, err := c.GetToken()
		Expect(t).To(Equal("test-token"))
		Expect(err).NotTo(HaveOccurred())
	})

	it("returns errors for GetToken()", func() {
		c.Creator = badCreator
		cf.GetTokenReturns("", errors.New("test-error-2"))
		t, err := c.GetToken()
		Expect(t).To(BeZero())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("test-error"))
	})

	it("passes through to ListOrgsByQuery()", func() {
		cf.ListOrgsByQueryReturns([]cfclient.Org{}, nil)
		t, err := c.ListOrgsByQuery(nil)
		Expect(t).NotTo(BeZero())
		Expect(err).NotTo(HaveOccurred())
	})

	it("returns errors for ListOrgsByQuery()", func() {
		c.Creator = badCreator
		cf.GetTokenReturns("", errors.New("test-error-2"))
		t, err := c.ListOrgsByQuery(nil)
		Expect(t).To(BeZero())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("test-error"))
	})

	it("passes through to CreateOrg()", func() {
		cf.CreateOrgReturns(cfclient.Org{Guid: "test-guid"}, nil)
		t, err := c.CreateOrg(cfclient.OrgRequest{})
		Expect(t).NotTo(BeZero())
		Expect(err).NotTo(HaveOccurred())
	})

	it("returns errors for CreateOrg()", func() {
		c.Creator = badCreator
		cf.GetTokenReturns("", errors.New("test-error-2"))
		t, err := c.CreateOrg(cfclient.OrgRequest{})
		Expect(t).To(BeZero())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("test-error"))
	})

	it("passes through to AssociateOrgUser()", func() {
		cf.AssociateOrgUserReturns(cfclient.Org{Guid: "test-guid"}, nil)
		t, err := c.AssociateOrgUser("org", "user")
		Expect(t).NotTo(BeZero())
		Expect(err).NotTo(HaveOccurred())
	})

	it("returns errors for AssociateOrgUser()", func() {
		c.Creator = badCreator
		cf.GetTokenReturns("", errors.New("test-error-2"))
		t, err := c.AssociateOrgUser("org", "user")
		Expect(t).To(BeZero())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("test-error"))
	})

	it("passes through to AssociateOrgAuditor()", func() {
		cf.AssociateOrgAuditorReturns(cfclient.Org{Guid: "test-guid"}, nil)
		t, err := c.AssociateOrgAuditor("org", "user")
		Expect(t).NotTo(BeZero())
		Expect(err).NotTo(HaveOccurred())
	})

	it("returns errors for AssociateOrgAuditor()", func() {
		c.Creator = badCreator
		cf.GetTokenReturns("", errors.New("test-error-2"))
		t, err := c.AssociateOrgAuditor("org", "user")
		Expect(t).To(BeZero())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("test-error"))
	})

	it("passes through to AssociateOrgManager()", func() {
		cf.AssociateOrgManagerReturns(cfclient.Org{Guid: "test-guid"}, nil)
		t, err := c.AssociateOrgManager("org", "user")
		Expect(t).NotTo(BeZero())
		Expect(err).NotTo(HaveOccurred())
	})

	it("returns errors for AssociateOrgManager()", func() {
		c.Creator = badCreator
		cf.GetTokenReturns("", errors.New("test-error-2"))
		t, err := c.AssociateOrgManager("org", "user")
		Expect(t).To(BeZero())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("test-error"))
	})

	it("passes through to GetOrgQuotaByName()", func() {
		cf.GetOrgQuotaByNameReturns(cfclient.OrgQuota{Guid: "test-guid"}, nil)
		t, err := c.GetOrgQuotaByName("quota")
		Expect(t).NotTo(BeZero())
		Expect(err).NotTo(HaveOccurred())
	})

	it("returns errors for GetOrgQuotaByName()", func() {
		c.Creator = badCreator
		cf.GetTokenReturns("", errors.New("test-error-2"))
		t, err := c.GetOrgQuotaByName("quota")
		Expect(t).To(BeZero())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("test-error"))
	})

	it("passes through to CreateSpace()", func() {
		cf.CreateSpaceReturns(cfclient.Space{Guid: "test-guid"}, nil)
		t, err := c.CreateSpace(cfclient.SpaceRequest{})
		Expect(t).NotTo(BeZero())
		Expect(err).NotTo(HaveOccurred())
	})

	it("returns errors for CreateSpace()", func() {
		c.Creator = badCreator
		cf.GetTokenReturns("", errors.New("test-error-2"))
		t, err := c.CreateSpace(cfclient.SpaceRequest{})
		Expect(t).To(BeZero())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("test-error"))
	})
}
