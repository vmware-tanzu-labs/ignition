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

func TestQuotaIDForName(t *testing.T) {
	spec.Run(t, "QuotaIDForName", testQuotaIDForName, spec.Report(report.Terminal{}))
}

func testQuotaIDForName(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	it("errors when the call to cloud foundry fails", func() {
		f := &cloudfoundryfakes.FakeAPI{}
		f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{}, errors.New("test error"))
		id, err := cloudfoundry.QuotaIDForName("test", f)
		Expect(id).To(BeZero())
		Expect(err).To(HaveOccurred())
	})

	it("returns the quota id when the call to cloud foundry succeeds", func() {
		f := &cloudfoundryfakes.FakeAPI{}
		f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{
			Guid: "test-org-quota-id",
		}, nil)
		id, err := cloudfoundry.QuotaIDForName("test", f)
		Expect(id).To(Equal("test-org-quota-id"))
		Expect(err).NotTo(HaveOccurred())
	})

	it("errors when the quota id is empty", func() {
		f := &cloudfoundryfakes.FakeAPI{}
		f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{}, nil)
		id, err := cloudfoundry.QuotaIDForName("test", f)
		Expect(id).To(BeZero())
		Expect(err).To(HaveOccurred())
	})
}
