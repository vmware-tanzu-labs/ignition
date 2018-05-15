package uaa_test

import (
	"net/http/httptest"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/internal"
	"github.com/pivotalservices/ignition/uaa"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestAuthenticate(t *testing.T) {
	spec.Run(t, "Authenticate", testAuthenticate, spec.Report(report.Terminal{}))
}

func testAuthenticate(t *testing.T, when spec.G, it spec.S) {
	var a *uaa.Client

	it.Before(func() {
		RegisterTestingT(t)
	})

	when("there is a need to authenticate but the token is empty", func() {
		var (
			s      *httptest.Server
			called bool
		)

		it.Before(func() {
			called = false
			calledFunc := func() {
				called = true
			}
			s = internal.ServeFromTestdata(t, "empty-token.json", calledFunc)
			a = &uaa.Client{
				URL:          s.URL,
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
			}
		})

		it.After(func() {
			s.Close()
		})

		it("returns an error", func() {
			err := a.Authenticate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("uaa: could not refresh token"))
			Expect(called).To(BeTrue())
		})
	})

	when("there is a need to authenticate and the token is valid", func() {
		var (
			s      *httptest.Server
			called bool
		)

		it.Before(func() {
			called = false
			calledFunc := func() {
				called = true
			}
			s = internal.ServeFromTestdata(t, "token.json", calledFunc)
			a = &uaa.Client{
				URL:          s.URL,
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
			}
		})

		it.After(func() {
			s.Close()
		})

		it("succeeds", func() {
			err := a.Authenticate()
			Expect(err).NotTo(HaveOccurred())
			Expect(called).To(BeTrue())
		})
	})
}
