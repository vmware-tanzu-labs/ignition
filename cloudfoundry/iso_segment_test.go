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

func TestISOSegmentForName(t *testing.T) {
	spec.Run(t, "ISOSegmentForName", testISOSegmentForName, spec.Report(report.Terminal{}))
}

func testISOSegmentForName(t *testing.T, when spec.G, it spec.S) {
	var f *cloudfoundryfakes.FakeAPI

	it.Before(func() {
		RegisterTestingT(t)
	})

	when("list iso segments returns an error", func() {
		it.Before(func() {
			f = &cloudfoundryfakes.FakeAPI{}
			f.ListIsolationSegmentsByQueryReturns(nil, errors.New("test error"))
		})

		it("errors when the call to cloud foundry fails", func() {
			id, err := cloudfoundry.ISOSegmentIDForName("test", f)
			Expect(id).To(BeZero())
			Expect(err).To(HaveOccurred())
		})
	})

	when("list iso segment only returns the shared default iso segment", func() {
		it.Before(func() {
			f = &cloudfoundryfakes.FakeAPI{}
			f.ListIsolationSegmentsByQueryReturns([]cfclient.IsolationSegment{
				cfclient.IsolationSegment{
					GUID: "shared-iso-guid",
					Name: "shared",
				},
			}, nil)
		})

		it("returns the shared iso segment name", func() {
			id, err := cloudfoundry.ISOSegmentIDForName("shared", f)
			Expect(err).ToNot(HaveOccurred())
			Expect(id).To(Equal("shared-iso-guid"))
		})

		it("returns err when not found", func() {
			id, err := cloudfoundry.ISOSegmentIDForName("doesnotexist", f)
			Expect(err).To(HaveOccurred())
			Expect(id).To(Equal(""))
		})
	})

	when("list iso segment returns the shared default iso segment and one other", func() {
		it.Before(func() {
			f = &cloudfoundryfakes.FakeAPI{}
			f.ListIsolationSegmentsByQueryReturns([]cfclient.IsolationSegment{
				cfclient.IsolationSegment{
					GUID: "shared-iso-guid",
					Name: "shared",
				},
				cfclient.IsolationSegment{
					GUID: "my-iso-guid",
					Name: "myiso",
				},
			}, nil)
		})

		it("returns the myiso segment name", func() {
			id, err := cloudfoundry.ISOSegmentIDForName("myiso", f)
			Expect(err).ToNot(HaveOccurred())
			Expect(id).To(Equal("my-iso-guid"))
		})

		it("returns err when not found", func() {
			id, err := cloudfoundry.ISOSegmentIDForName("doesnotexist", f)
			Expect(err).To(HaveOccurred())
			Expect(id).To(Equal(""))
		})

		it("returns err when name is empty", func() {
			id, err := cloudfoundry.ISOSegmentIDForName("", f)
			Expect(err).To(HaveOccurred())
			Expect(id).To(Equal(""))
		})
	})
}
