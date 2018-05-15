package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pivotalservices/ignition/cloudfoundry"
)

// Info is metadata that ignition API clients can use to display their UX
type Info struct {
	CompanyName              string
	ExperimentationSpaceName string
	IgnitionOrgCount         int
}

// InfoHandler writes the contents of the provided Info to the response
func InfoHandler(
	companyName, spaceName, orgQuotaID string,
	updateFreq time.Duration,
	orgQuerier cloudfoundry.OrganizationQuerier) http.Handler {

	orgCount := getIgnitionOrgCount(orgQuotaID, orgQuerier)
	startBackgroundOrgCountUpdater(orgQuotaID, orgQuerier, &orgCount, updateFreq)

	fn := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		i := Info{
			CompanyName:              companyName,
			ExperimentationSpaceName: spaceName,
			IgnitionOrgCount:         orgCountOrDefault(orgCount),
		}
		json.NewEncoder(w).Encode(i)
	}
	return http.HandlerFunc(fn)
}

func getIgnitionOrgCount(orgQuotaID string, orgQuerier cloudfoundry.OrganizationQuerier) *int {
	orgCount, err := queryIgnitionOrgCount(orgQuotaID, orgQuerier)
	if err != nil {
		// ignition org count is non-critical - so log it and continue
		log.Println(fmt.Sprintf("[ERROR] Could not get updated org count: %s", err.Error()))
	}
	return orgCount
}

func queryIgnitionOrgCount(orgQuotaID string, orgQuerier cloudfoundry.OrganizationQuerier) (*int, error) {
	orgs, err := orgQuerier.ListOrgsByQuery(url.Values{})
	if err != nil {
		return nil, err
	}

	count := 0
	for _, o := range orgs {
		if o.QuotaDefinitionGuid == orgQuotaID {
			count++
		}
	}
	return &count, nil
}

func startBackgroundOrgCountUpdater(
	orgQuotaID string,
	orgQuerier cloudfoundry.OrganizationQuerier,
	orgCount **int,
	updateFreq time.Duration) {
	go func() {
		for {
			time.Sleep(updateFreq)
			oc := getIgnitionOrgCount(orgQuotaID, orgQuerier)
			if oc != nil && *oc > 0 {
				*(orgCount) = oc
			}
		}
	}()
}

func orgCountOrDefault(orgCount *int) int {
	o := 0
	if orgCount != nil {
		o = *orgCount
	}
	return o
}
