package cloudfoundry

import (
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pkg/errors"
)

// SpaceCreator creates spaces
type SpaceCreator interface {
	CreateSpace(req cfclient.SpaceRequest) (cfclient.Space, error)
}

// CreateSpace creates a space with the given name
// for the given user
func CreateSpace(name string, organizationID string, userID string, a SpaceCreator) error {
	req := cfclient.SpaceRequest{
		Name:             strings.ToLower(name),
		AuditorGuid:      []string{userID},
		DeveloperGuid:    []string{userID},
		ManagerGuid:      []string{userID},
		OrganizationGuid: organizationID,
		AllowSSH:         true,
	}
	space, err := a.CreateSpace(req)
	if err != nil || space.Guid == "" {
		return errors.Wrapf(err, "could not create space with name [%s] and organizationID [%s]", name, organizationID)
	}

	return nil
}
