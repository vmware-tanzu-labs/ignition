package uaa

import (
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/uaa-cli/uaa"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Client provides access to the UAA API
type Client struct {
	URL          string
	ClientID     string
	ClientSecret string
	Client       *http.Client
	Config       *clientcredentials.Config
	Token        *oauth2.Token
	userManager  *uaa.UserManager
}

// UserIDForAccountName queries the UAA API for users filtered by account name
func (a *Client) UserIDForAccountName(accountName string) (string, error) {
	if strings.TrimSpace(accountName) == "" {
		return "", errors.New("cannot search for a user with an empty account name")
	}

	err := a.Authenticate()
	if err != nil {
		return "", errors.Wrap(err, "uaa: cannot authenticate")
	}
	user, err := a.userManager.GetByUsername(accountName, "", "")
	if err != nil {
		return "", errors.Wrap(err, "uaa: cannot get user")
	}
	if strings.TrimSpace(user.ID) == "" {
		return "", errors.Errorf("cannot find user with account name: [%s]", accountName)
	}
	return user.ID, nil
}

// CreateUser creates new users in the UAA database.
func (a *Client) CreateUser(username, origin, externalID, email string) (string, error) {
	err := a.Authenticate()
	if err != nil {
		return "", errors.Wrap(err, "uaa: cannot authenticate")
	}

	user, err := a.userManager.Create(uaa.ScimUser{
		Username:   username,
		Origin:     origin,
		ExternalId: externalID,
		Emails:     []uaa.ScimUserEmail{{Value: email}},
	})

	if err != nil {
		return "", errors.Wrap(err, "uaa: cannot create user "+username)
	}
	return user.ID, nil
}
