package cloudfoundry

import (
	"fmt"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

// QuotaQuerier is used to query a Cloud Controller API for quotas
type QuotaQuerier interface {
	GetOrgQuotaByName(name string) (cfclient.OrgQuota, error)
}

// QuotaIDForName gets the quota ID for the given quota name
func QuotaIDForName(name string, q QuotaQuerier) (string, error) {
	quota, err := q.GetOrgQuotaByName(name)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(quota.Guid) == "" {
		return "", fmt.Errorf("cannot find quota with name [%s]", name)
	}
	return quota.Guid, nil
}
