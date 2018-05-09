package cloudfoundry

import (
	"fmt"
	"net/url"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pkg/errors"
)

// Organization is a Cloud Foundry Organization
type Organization struct {
	GUID                        string `json:"guid"`
	CreatedAt                   string `json:"created_at"`
	UpdatedAt                   string `json:"updated_at"`
	Name                        string `json:"name"`
	QuotaDefinitionGUID         string `json:"quota_definition_guid"`
	DefaultIsolationSegmentGUID string `json:"default_isolation_segment_guid"`
	URL                         string `json:"url"`
}

// OrganizationQuerier is used to query a Cloud Controller API for organizations
type OrganizationQuerier interface {
	ListOrgsByQuery(query url.Values) ([]cfclient.Org, error)
}

// OrganizationCreator creates orgs
type OrganizationCreator interface {
	CreateOrg(req cfclient.OrgRequest) (cfclient.Org, error)
	UpdateOrg(orgGUID string, orgRequest cfclient.OrgRequest) (cfclient.Org, error)
	AddIsolationSegmentToOrg(isolationSegmentGUID, orgGUID string) error
}

// RoleGrantor allows for users to be granted org and space roles
type RoleGrantor interface {
	AssociateOrgUser(orgGUID, userGUID string) (cfclient.Org, error)
	AssociateOrgAuditor(orgGUID, userGUID string) (cfclient.Org, error)
	AssociateOrgManager(orgGUID, userGUID string) (cfclient.Org, error)
}

// OrgsForUserID returns the orgs that the user is a member of
func OrgsForUserID(id string, appsURL string, q OrganizationQuerier) ([]Organization, error) {
	query := url.Values{}
	query.Add("q", fmt.Sprintf("user_guid:%s", id))
	o, err := q.ListOrgsByQuery(query)
	if err != nil {
		return nil, err
	}

	result := make([]Organization, len(o))
	for i := range o {
		result[i] = convertOrg(o[i], appsURL)
	}
	return result, nil
}

// CreateOrg creates an organization with the given name and quota for
// the given user
func CreateOrg(name, appsURL, quotaID, isoSegmentID string, a OrganizationCreator) (*Organization, error) {
	req := cfclient.OrgRequest{
		Name:                strings.ToLower(name),
		QuotaDefinitionGuid: quotaID,
	}
	org, err := a.CreateOrg(req)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create org with name [%s] and quota [%s]", name, quotaID)
	}

	// enable the org to use the iso segment
	err = a.AddIsolationSegmentToOrg(isoSegmentID, org.Guid)
	if err != nil {
		return nil, errors.Wrapf(err, "could not assign isolation segment [%s] to org [%s]", isoSegmentID, name)
	}

	// make the iso segment the default for the org
	req.DefaultIsolationSegmentGuid = isoSegmentID
	_, err = a.UpdateOrg(org.Guid, req)
	if err != nil {
		return nil, errors.Wrapf(err, "could not make isolation segment [%s] the default for org [%s]", isoSegmentID, name)
	}

	o := convertOrg(org, appsURL)
	return &o, nil
}

func convertOrg(o cfclient.Org, appsURL string) Organization {
	return Organization{
		GUID:                        o.Guid,
		CreatedAt:                   o.CreatedAt,
		UpdatedAt:                   o.UpdatedAt,
		Name:                        o.Name,
		QuotaDefinitionGUID:         o.QuotaDefinitionGuid,
		DefaultIsolationSegmentGUID: o.DefaultIsolationSegmentGuid,
		URL: fmt.Sprintf("%s/organizations/%s", appsURL, o.Guid),
	}
}
