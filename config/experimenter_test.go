package config

import (
	"errors"
	"os"
	"testing"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/cloudfoundry/cloudfoundryfakes"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestNewExperimenter(t *testing.T) {
	spec.Run(t, "NewExperimenter", testNewExperimenter, spec.Report(report.Terminal{}))
}

func testNewExperimenter(t *testing.T, when spec.G, it spec.S) {
	var f *cloudfoundryfakes.FakeAPI
	reset := func() {
		os.Unsetenv("VCAP_APPLICATION")
		os.Unsetenv("VCAP_SERVICES")
		os.Unsetenv("PORT")

		os.Unsetenv("IGNITION_ORG_PREFIX")
		os.Unsetenv("IGNITION_QUOTA_NAME")
		os.Unsetenv("IGNITION_SPACE_NAME")
	}

	it.Before(func() {
		RegisterTestingT(t)
		f = &cloudfoundryfakes.FakeAPI{}
		reset()
	})

	it.After(func() {
		reset()
	})

	when("not running on CF", func() {
		it("succeeds when no variables are set", func() {
			f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{
				Guid: "test-quota-id",
			}, nil)
			e, err := NewExperimenter("ignition-config", f)
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())
			Expect(e.OrgPrefix).To(Equal("ignition"))
			Expect(e.QuotaName).To(Equal("ignition"))
			Expect(e.QuotaID).NotTo(BeZero())
			Expect(e.SpaceName).To(Equal("playground"))
		})

		it("looks up the Quota ID if it is missing", func() {
			f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{
				Guid: "test-quota-id",
			}, nil)
			e, err := NewExperimenter("ignition-config", f)
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())
			Expect(e.OrgPrefix).To(Equal("ignition"))
			Expect(e.QuotaName).To(Equal("ignition"))
			Expect(e.QuotaID).To(Equal("test-quota-id"))
			Expect(e.SpaceName).To(Equal("playground"))
		})

		it("falls back to the default quota if the quota cannot be found", func() {
			f.GetOrgQuotaByNameReturnsOnCall(0, cfclient.OrgQuota{}, errors.New("not found"))
			f.GetOrgQuotaByNameReturnsOnCall(1, cfclient.OrgQuota{
				Guid: "default-quota-id",
			}, nil)
			e, err := NewExperimenter("ignition-config", f)
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())
			Expect(e.OrgPrefix).To(Equal("ignition"))
			Expect(e.QuotaName).To(Equal("default"))
			Expect(e.QuotaID).To(Equal("default-quota-id"))
			Expect(e.SpaceName).To(Equal("playground"))
		})

		it("errors if the named and the default quota cannot be found", func() {
			f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{}, errors.New("not found"))
			e, err := NewExperimenter("ignition-config", f)
			Expect(err).To(HaveOccurred())
			Expect(e).To(BeNil())
		})

		when("the quota name is set but empty", func() {
			it.Before(func() {
				os.Setenv("IGNITION_QUOTA_NAME", "   ")
			})

			it.After(func() {
				os.Unsetenv("IGNITION_QUOTA_NAME")
			})

			it("uses the default quota", func() {
				f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{
					Guid: "default-quota-id",
				}, nil)
				e, err := NewExperimenter("ignition-config", f)
				Expect(err).NotTo(HaveOccurred())
				Expect(e).NotTo(BeNil())
				Expect(e.OrgPrefix).To(Equal("ignition"))
				Expect(e.QuotaName).To(Equal("default"))
				Expect(e.QuotaID).To(Equal("default-quota-id"))
				Expect(e.SpaceName).To(Equal("playground"))
			})
		})
	})

	when("running on CF", func() {
		it.Before(func() {
			os.Setenv("VCAP_APPLICATION", `{"cf_api": "https://api.run.pcfbeta.io","limits": {"fds": 16384},"application_name": "ignition","application_uris": ["ignition.pcfbeta.io"],"name": "ignition","space_name": "development","space_id": "test-space-id","uris": ["ignition.pcfbeta.io"],"users": null,"application_id": "test-app-id"}`)
			os.Setenv("VCAP_SERVICES", `{}`)
			os.Setenv("PORT", "54321")
		})

		it("errors when VCAP_APPLICATION contents are invalid", func() {
			os.Setenv("VCAP_APPLICATION", "%&^%@")
			e, err := NewExperimenter("ignition-config", f)
			Expect(err).To(HaveOccurred())
			Expect(e).To(BeNil())
		})

		it("succeeds when no config exists", func() {
			f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{
				Guid: "test-quota-id",
			}, nil)
			e, err := NewExperimenter("ignition-config", f)
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())
			Expect(e.OrgPrefix).To(Equal("ignition"))
			Expect(e.QuotaName).To(Equal("ignition"))
			Expect(e.QuotaID).To(Equal("test-quota-id"))
			Expect(e.SpaceName).To(Equal("playground"))
		})

		it("uses an org prefix specified in ignition-config", func() {
			os.Setenv("VCAP_SERVICES", `{"user-provided": [{
				"name": "ignition-config",
				"instance_name": "ignition-config",
				"credentials": {
					"org_prefix": "test-org-prefix"
				}}]}`)
			f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{
				Guid: "test-quota-id",
			}, nil)
			e, err := NewExperimenter("ignition-config", f)
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())
			Expect(e.OrgPrefix).To(Equal("test-org-prefix"))
			Expect(e.QuotaName).To(Equal("ignition"))
			Expect(e.QuotaID).NotTo(BeZero())
			Expect(e.SpaceName).To(Equal("playground"))
		})

		it("uses a quota name specified in ignition-config", func() {
			os.Setenv("VCAP_SERVICES", `{"user-provided": [{
				"name": "ignition-config",
				"instance_name": "ignition-config",
				"credentials": {
					"quota_name": "test-ignition-quota-name"
				}}]}`)
			f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{
				Guid: "test-quota-id",
			}, nil)
			e, err := NewExperimenter("ignition-config", f)
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())
			Expect(e.OrgPrefix).To(Equal("ignition"))
			Expect(e.QuotaName).To(Equal("test-ignition-quota-name"))
			Expect(e.QuotaID).NotTo(BeZero())
			Expect(e.SpaceName).To(Equal("playground"))
		})

		it("uses a space name specified in ignition-config", func() {
			os.Setenv("VCAP_SERVICES", `{"user-provided": [{
				"name": "ignition-config",
				"instance_name": "ignition-config",
				"credentials": {
					"space_name": "test-ignition-space-name"
				}}]}`)
			f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{
				Guid: "test-space-id",
			}, nil)
			e, err := NewExperimenter("ignition-config", f)
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())
			Expect(e.OrgPrefix).To(Equal("ignition"))
			Expect(e.QuotaName).To(Equal("ignition"))
			Expect(e.QuotaID).NotTo(BeZero())
			Expect(e.SpaceName).To(Equal("test-ignition-space-name"))
		})
	})
}
