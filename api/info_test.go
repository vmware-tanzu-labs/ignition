package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pivotalservices/ignition/api"
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

	it("", func() {
		i := api.Info{
			CompanyName:              "Test Company",
			ExperimentationSpaceName: "Test Space",
		}
		handler := api.InfoHandler(i)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
		Expect(w.Code).To(Equal(http.StatusOK))
		Expect(w.Body.String()).To(ContainSubstring("Test Company"))
		Expect(w.Body.String()).To(ContainSubstring("Test Space"))
	})
}
