package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pivotalservices/ignition/user"
	"github.com/pivotalservices/ignition/user/openid"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Authorizer is used to authenticate and authorize the user
type Authorizer struct {
	Variant      string          `envconfig:"auth_variant" default:"p-identity"`                    // IGNITION_AUTH_VARIANT
	ServiceName  string          `envconfig:"auth_servicename" default:"ignition-identity"`         // IGNITION_AUTH_SERVICENAME
	ClientID     string          `envconfig:"client_id"`                                            // IGNITION_CLIENT_ID << REQUIRED
	ClientSecret string          `envconfig:"client_secret"`                                        // IGNITION_CLIENT_SECRET << REQUIRED
	URL          string          `envconfig:"auth_url"`                                             // IGNITION_AUTH_URL << REQUIRED
	Domain       string          `envconfig:"authorized_domain"`                                    // IGNITION_AUTHORIZED_DOMAIN << REQUIRED
	Scopes       []string        `envconfig:"auth_scopes" default:"openid,profile,user_attributes"` // IGNITION_AUTH_SCOPES
	Provider     *Provider       `ignored:"true"`
	Verifier     openid.Verifier `ignored:"true"`
	Fetcher      user.Fetcher    `ignored:"true"`
	Config       *oauth2.Config  `ignored:"true"`
}

// Provider is an OpenID Connect provider
type Provider struct {
	Issuer          string   `json:"issuer"`
	AuthURL         string   `json:"authorization_endpoint"`
	TokenURL        string   `json:"token_endpoint"`
	JWKSURL         string   `json:"jwks_uri"`
	UserInfoURL     string   `json:"userinfo_endpoint"`
	ScopesSupported []string `json:"scopes_supported"`
}

// NewAuthorizer uses environment variables to populate a new Authorizer
func NewAuthorizer(name string) (*Authorizer, error) {
	var a Authorizer
	envconfig.Process(ignition, &a)
	if cfenv.IsRunningOnCF() {
		c, err := cfenv.Current()
		if err != nil {
			return nil, err
		}
		s, err := c.Services.WithName(name)
		if err != nil {
			return nil, err
		}

		variant, ok := s.CredentialString("auth_variant")
		if ok && strings.TrimSpace(variant) != "" {
			a.Variant = variant
		}

		domain, ok := s.CredentialString("authorized_domain")
		if ok && strings.TrimSpace(domain) != "" {
			a.Domain = domain
		}

		authURL, ok := s.CredentialString("auth_url")
		if ok && strings.TrimSpace(authURL) != "" {
			a.URL = authURL
		}

		clientID, ok := s.CredentialString("client_id")
		if ok && strings.TrimSpace(clientID) != "" {
			a.ClientID = clientID
		}
		clientSecret, ok := s.CredentialString("client_secret")
		if ok && strings.TrimSpace(clientSecret) != "" {
			a.ClientSecret = clientSecret
		}
		serviceName, ok := s.CredentialString("auth_servicename")
		if ok && strings.TrimSpace(serviceName) != "" {
			a.ServiceName = serviceName
		}

		if strings.EqualFold(strings.TrimSpace(a.Variant), "p-identity") {
			i, err := c.Services.WithName(a.ServiceName)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("a Single Sign On service instance with the name \"%s\" is required to use this app", a.ServiceName))
			}

			authURL, ok := i.CredentialString("auth_domain")
			if !ok {
				return nil, fmt.Errorf("could not retrieve the client_id; make sure you have created and bound a Single Sign On service instance with the name \"%s\"", a.ServiceName)
			}
			a.URL = authURL
			clientid, ok := i.CredentialString("client_id")
			if !ok {
				return nil, fmt.Errorf("could not retrieve the client_id; make sure you have created and bound a Single Sign On service instance with the name \"%s\"", a.ServiceName)
			}
			a.ClientID = clientid
			clientsecret, ok := i.CredentialString("client_secret")
			if !ok {
				return nil, fmt.Errorf("could not retrieve the client_secret; make sure you have created and bound a Single Sign On service instance with the name \"%s\"", a.ServiceName)
			}
			a.ClientSecret = clientsecret
		}
	}

	if strings.TrimSpace(a.ClientID) == "" {
		return nil, errors.New("client_id is required")
	}
	if strings.TrimSpace(a.ClientSecret) == "" {
		return nil, errors.New("client_secret is required")
	}
	if strings.TrimSpace(a.URL) == "" {
		return nil, errors.New("auth_url is required")
	}
	if strings.TrimSpace(a.Domain) == "" {
		return nil, errors.New("authorized_domain is required")
	}

	wellKnown := strings.TrimSuffix(a.URL, "/") + "/.well-known/openid-configuration"
	resp, err := http.Get(wellKnown)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, body)
	}

	var p Provider
	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, err
	}
	a.Provider = &p
	// TODO: Warn when a.Scopes includes items that are not in p.ScopesSupported
	a.Verifier = openid.NewVerifier(p.Issuer, a.ClientID, p.JWKSURL)
	a.Fetcher = &openid.Fetcher{
		Verifier: a.Verifier,
	}
	a.Config = &oauth2.Config{
		ClientID:     a.ClientID,
		ClientSecret: a.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  p.AuthURL,
			TokenURL: p.TokenURL,
		},
		Scopes: a.Scopes,
	}
	return &a, nil
}
