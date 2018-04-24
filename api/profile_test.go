package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/api"
	"github.com/pivotalservices/ignition/user"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestProfileHandler(t *testing.T) {
	spec.Run(t, "ProfileHandler", testProfileHandler, spec.Report(report.Terminal{}))
}

func testProfileHandler(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})

	it("returns unauthorized when there is no profile", func() {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		api.ProfileHandler().ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusUnauthorized))
	})

	it("writes the profile when it is in the context when from a different domain", func() {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(user.WithProfile(req.Context(), &user.Profile{
			AccountName: "test@example.net",
			Name:        "Test User",
			Email:       "test@example.net",
		}))
		w := httptest.NewRecorder()
		api.ProfileHandler().ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.String()).To(ContainSubstring("test@example.net"))
	})

	it("writes the profile when it is in the context", func() {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(user.WithProfile(req.Context(), &user.Profile{
			AccountName: "test@example.net",
			Name:        "Test User",
			Email:       "test@example.net",
		}))
		w := httptest.NewRecorder()
		api.ProfileHandler().ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.String()).To(ContainSubstring("test@example.net"))
	})
}
