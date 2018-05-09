package config

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

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
		os.Unsetenv("IGNITION_ORG_COUNT_UPDATE_INTERVAL")
		os.Unsetenv("IGNITION_QUOTA_NAME")
		os.Unsetenv("IGNITION_SPACE_NAME")
		os.Unsetenv("IGNITION_ISO_SEGMENT_NAME")
	}

	it.Before(func() {
		RegisterTestingT(t)
		reset()
		f = &cloudfoundryfakes.FakeAPI{}
		f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{
			Guid: "test-quota-id",
		}, nil)
		f.ListIsolationSegmentsByQueryReturns([]cfclient.IsolationSegment{
			cfclient.IsolationSegment{
				Name: "shared",
				GUID: "shared-iso-segment-id",
			},
		}, nil)
	})

	it.After(func() {
		reset()
	})

	when("not running on CF", func() {
		it("succeeds when no variables are set", func() {
			e := createExperimenter(f)
			Expect(e.OrgPrefix).To(Equal("ignition"))
			Expect(e.OrgCountUpdateInterval).To(Equal(time.Minute))
			Expect(e.QuotaName).To(Equal("ignition"))
			Expect(e.QuotaID).NotTo(BeZero())
			Expect(e.SpaceName).To(Equal("playground"))
			Expect(e.ISOSegmentName).To(Equal("shared"))
			Expect(e.ISOSegmentID).NotTo(BeZero())
		})

		it("looks up the Quota ID if it is missing", func() {
			e := createExperimenter(f)
			Expect(e.QuotaID).To(Equal("test-quota-id"))
		})

		it("looks up the ISO Segment ID if it is missing", func() {
			e := createExperimenter(f)
			Expect(e.ISOSegmentID).To(Equal("shared-iso-segment-id"))
		})

		it("falls back to the default quota if the quota cannot be found", func() {
			f.GetOrgQuotaByNameReturnsOnCall(0, cfclient.OrgQuota{}, errors.New("not found"))
			f.GetOrgQuotaByNameReturnsOnCall(1, cfclient.OrgQuota{
				Guid: "default-quota-id",
			}, nil)
			e := createExperimenter(f)
			Expect(e.OrgPrefix).To(Equal("ignition"))
			Expect(e.OrgCountUpdateInterval).To(Equal(time.Minute))
			Expect(e.QuotaName).To(Equal("default"))
			Expect(e.QuotaID).To(Equal("default-quota-id"))
			Expect(e.SpaceName).To(Equal("playground"))
		})

		it("errors if the named and the default quota cannot be found", func() {
			f.GetOrgQuotaByNameReturns(cfclient.OrgQuota{}, errors.New("not found"))
			e, err := NewExperimenter("ignition-config", f, f)
			Expect(err).To(HaveOccurred())
			Expect(e).To(BeNil())
		})

		when("the ignition environment variables are set", func() {
			it.Before(func() {
				os.Setenv("IGNITION_ORG_PREFIX", "env-org")
				os.Setenv("IGNITION_ORG_COUNT_UPDATE_INTERVAL", "5m")
				os.Setenv("IGNITION_SPACE_NAME", "env-space")
				os.Setenv("IGNITION_QUOTA_NAME", "env-quota-name")
				os.Setenv("IGNITION_ISO_SEGMENT_NAME", "env-iso-segment-name")
				f.ListIsolationSegmentsByQueryReturns([]cfclient.IsolationSegment{
					cfclient.IsolationSegment{
						Name: "env-iso-segment-name",
						GUID: "env-iso-segment-id",
					},
				}, nil)
			})

			it("uses the values specified in the individual environment variables", func() {
				e := createExperimenter(f)
				Expect(e.OrgPrefix).To(Equal("env-org"))
				Expect(e.OrgCountUpdateInterval).To(Equal(time.Minute * 5))
				Expect(e.QuotaName).To(Equal("env-quota-name"))
				Expect(e.QuotaID).To(Equal("test-quota-id"))
				Expect(e.SpaceName).To(Equal("env-space"))
				Expect(e.ISOSegmentName).To(Equal("env-iso-segment-name"))
				Expect(e.ISOSegmentID).To(Equal("env-iso-segment-id"))
			})
		})

		when("the quota name is set but empty", func() {
			it.Before(func() {
				os.Setenv("IGNITION_QUOTA_NAME", "   ")
			})

			it("uses the default quota", func() {
				e := createExperimenter(f)
				Expect(e.OrgPrefix).To(Equal("ignition"))
				Expect(e.QuotaName).To(Equal("default"))
				Expect(e.QuotaID).To(Equal("test-quota-id"))
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
			e, err := NewExperimenter("ignition-config", f, f)
			Expect(err).To(HaveOccurred())
			Expect(e).To(BeNil())
		})

		it("succeeds and uses the defaults when no config exists", func() {
			e := createExperimenter(f)
			Expect(e.OrgPrefix).To(Equal("ignition"))
			Expect(e.OrgCountUpdateInterval).To(Equal(time.Minute))
			Expect(e.QuotaName).To(Equal("ignition"))
			Expect(e.QuotaID).To(Equal("test-quota-id"))
			Expect(e.SpaceName).To(Equal("playground"))
		})

		it("uses the org prefix specified in ignition-config", func() {
			stubCupsService("org_prefix", "test-org-prefix")
			e := createExperimenter(f)
			Expect(e.OrgPrefix).To(Equal("test-org-prefix"))
		})

		it("uses the quota name specified in ignition-config", func() {
			stubCupsService("quota_name", "test-ignition-quota-name")
			e := createExperimenter(f)
			Expect(e.QuotaName).To(Equal("test-ignition-quota-name"))
			Expect(e.QuotaID).To(Equal("test-quota-id"))
		})

		it("uses the space name specified in ignition-config", func() {
			stubCupsService("space_name", "test-ignition-space-name")
			e := createExperimenter(f)
			Expect(e.SpaceName).To(Equal("test-ignition-space-name"))
		})

		it("uses the org count update interval specified in ignition-config", func() {
			stubCupsService("org_count_update_interval", "3m")
			e := createExperimenter(f)
			Expect(e.OrgCountUpdateInterval).To(Equal(time.Minute * 3))
		})

		it("defaults the org count update interval to 1m when given an invalid duration", func() {
			stubCupsService("org_count_update_interval", "garbage")
			e := createExperimenter(f)
			Expect(e.OrgCountUpdateInterval).To(Equal(time.Minute))
		})

		it("uses the isolation segment name specified in ignition-config", func() {
			stubCupsService("iso_segment_name", "test-ignition-iso-segment-name")
			f.ListIsolationSegmentsByQueryReturns([]cfclient.IsolationSegment{
				cfclient.IsolationSegment{
					Name: "test-ignition-iso-segment-name",
					GUID: "test-iso-segment-id",
				},
			}, nil)
			e := createExperimenter(f)
			Expect(e.ISOSegmentName).To(Equal("test-ignition-iso-segment-name"))
			Expect(e.ISOSegmentID).To(Equal("test-iso-segment-id"))
		})

		it("defaults the isolation segment name to shared", func() {
			stubCupsService("iso_segment_name", "")
			e := createExperimenter(f)
			Expect(e.ISOSegmentName).To(Equal("shared"))
			Expect(e.ISOSegmentID).To(Equal("shared-iso-segment-id"))
		})
	})
}

func createExperimenter(f *cloudfoundryfakes.FakeAPI) *Experimenter {
	e, err := NewExperimenter("ignition-config", f, f)
	Expect(err).NotTo(HaveOccurred())
	Expect(e).NotTo(BeNil())
	return e
}

func stubCupsService(key, value string) {
	os.Setenv("VCAP_SERVICES", fmt.Sprintf(`{"user-provided": [{
		"name": "ignition-config",
		"instance_name": "ignition-config",
		"credentials": {
			"%s": "%s"
		}}]}`, key, value))
}
