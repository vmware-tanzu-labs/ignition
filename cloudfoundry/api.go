package cloudfoundry

import (
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient"
)

// API is a Cloud Controller API
type API interface {
	OrganizationCreator
	OrganizationQuerier
	SpaceCreator
	RoleGrantor
	QuotaQuerier
	GetToken() (string, error)
}

// Client is a Cloud Foundry client that automatically refreshes the token used
// to access the Cloud Controller API
type Client struct {
	CF      API
	Creator func(config *cfclient.Config) (client API, err error)
	Config  *cfclient.Config
}

func (c *Client) checkAuthentication() error {
	if c.CF == nil {
		return c.authenticate()
	}

	_, err := c.CF.GetToken()
	if err != nil {
		return c.authenticate()
	}
	return nil
}

func (c *Client) authenticate() error {
	var (
		cf  API
		err error
	)
	if c.Creator == nil {
		cf, err = cfclient.NewClient(c.Config)
	} else {
		cf, err = c.Creator(c.Config)
	}

	if err != nil {
		return err
	}
	c.CF = cf
	return nil
}

// GetToken is used to validate that the current client has valid authentication
// to the Cloud Controller API
func (c *Client) GetToken() (string, error) {
	err := c.checkAuthentication()
	if err != nil {
		return "", err
	}
	return c.CF.GetToken()
}

// ListOrgsByQuery uses the given query to request a filtered list of orgs
func (c *Client) ListOrgsByQuery(query url.Values) ([]cfclient.Org, error) {
	err := c.checkAuthentication()
	if err != nil {
		return nil, err
	}
	return c.CF.ListOrgsByQuery(query)
}

// CreateOrg creates a Cloud Foundry organization with the given OrgRequest
func (c *Client) CreateOrg(req cfclient.OrgRequest) (cfclient.Org, error) {
	err := c.checkAuthentication()
	if err != nil {
		return cfclient.Org{}, err
	}
	return c.CF.CreateOrg(req)
}

// AssociateOrgUser grants the given user the org user role for the given
// org
func (c *Client) AssociateOrgUser(orgGUID, userGUID string) (cfclient.Org, error) {
	err := c.checkAuthentication()
	if err != nil {
		return cfclient.Org{}, err
	}
	return c.CF.AssociateOrgUser(orgGUID, userGUID)
}

// AssociateOrgAuditor grants the given user the org auditor role for the given
// org
func (c *Client) AssociateOrgAuditor(orgGUID, userGUID string) (cfclient.Org, error) {
	err := c.checkAuthentication()
	if err != nil {
		return cfclient.Org{}, err
	}
	return c.CF.AssociateOrgAuditor(orgGUID, userGUID)
}

// AssociateOrgManager grants the given user the org manager role for the given
// org
func (c *Client) AssociateOrgManager(orgGUID, userGUID string) (cfclient.Org, error) {
	err := c.checkAuthentication()
	if err != nil {
		return cfclient.Org{}, err
	}
	return c.CF.AssociateOrgManager(orgGUID, userGUID)
}

// GetOrgQuotaByName will return the org quota with the given name, if it exists
func (c *Client) GetOrgQuotaByName(name string) (cfclient.OrgQuota, error) {
	err := c.checkAuthentication()
	if err != nil {
		return cfclient.OrgQuota{}, err
	}
	return c.CF.GetOrgQuotaByName(name)
}

// CreateSpace will create a Cloud Foundry space with the given SpaceRequest
func (c *Client) CreateSpace(req cfclient.SpaceRequest) (cfclient.Space, error) {
	err := c.checkAuthentication()
	if err != nil {
		return cfclient.Space{}, err
	}
	return c.CF.CreateSpace(req)
}
