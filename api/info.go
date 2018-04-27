package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/pivotalservices/ignition/cloudfoundry"
)

// Info is metadata that ignition API clients can use to display their UX
type Info struct {
	CompanyName              string
	ExperimentationSpaceName string
	IgnitionOrgCount         int
}

// InfoHandler writes the contents of the provided Info to the response
func InfoHandler(companyName, spaceName, orgQuotaID string, orgQuerier cloudfoundry.OrganizationQuerier) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		orgCount, err := ignitionOrgCount(orgQuotaID, orgQuerier)
		if err != nil {
			// ignition org count is non-critical - so log it and continue
			log.Println(err)
		}

		i := Info{
			CompanyName:              companyName,
			ExperimentationSpaceName: spaceName,
			IgnitionOrgCount:         orgCount,
		}
		json.NewEncoder(w).Encode(i)
	}
	return http.HandlerFunc(fn)
}

func ignitionOrgCount(orgQuotaID string, orgQuerier cloudfoundry.OrganizationQuerier) (int, error) {
	orgs, err := orgQuerier.ListOrgsByQuery(url.Values{})
	if err != nil {
		return 0, err
	}

	count := 0
	for _, o := range orgs {
		if o.QuotaDefinitionGuid == orgQuotaID {
			count++
		}
	}
	return count, nil
}
