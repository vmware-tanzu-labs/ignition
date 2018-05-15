package uaa

import (
	"context"
	"fmt"

	"github.com/cloudfoundry-incubator/uaa-cli/uaa"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/clientcredentials"
)

// API is used to access a UAA server
type API interface {
	UserIDForAccountName(a string) (string, error)
	CreateUser(username, origin, externalID, email string) (string, error)
}

// Authenticate will authenticate with a UAA server and set the Token and Client
// for the UAAAPI
func (a *Client) Authenticate() error {
	if a.Config == nil {
		a.Config = &clientcredentials.Config{
			ClientID:     a.ClientID,
			ClientSecret: a.ClientSecret,
			TokenURL:     fmt.Sprintf("%s/oauth/token", a.URL),
			Scopes:       []string{"cloud_controller.admin", "scim.write", "scim.read"},
		}
	}

	if a.Client == nil {
		a.Client = a.Config.Client(context.Background())
	}

	newToken := false
	if a.Token == nil || !a.Token.Valid() {
		token, err := a.Config.Token(context.Background())
		if err != nil {
			return errors.Wrap(err, "uaa: could not refresh token")
		}
		a.Token = token
		newToken = true
	}

	if newToken || a.userManager == nil {
		uaaConfig := uaa.NewConfig()
		uaaConfig.AddTarget(uaa.Target{BaseUrl: a.URL})
		uaaConfig.AddContext(uaa.NewContextWithToken(a.Token.AccessToken))
		a.userManager = &uaa.UserManager{
			Config:     uaaConfig,
			HttpClient: a.Client,
		}
	}

	return nil
}
