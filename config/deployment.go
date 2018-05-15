package config

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pivotalservices/ignition/cloudfoundry"
	"github.com/pivotalservices/ignition/uaa"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/clientcredentials"
)

// Deployment is a Cloud Foundry Deployment
type Deployment struct {
	SystemDomain      string           `envconfig:"system_domain"`                       // IGNITION_SYSTEM_DOMAIN << REQUIRED
	AppsURL           string           `ignored:"true"`                                  // Ignored
	APIURL            string           `ignored:"true"`                                  // Ignored
	UAAURL            string           `ignored:"true"`                                  // Ignored
	UAAOrigin         string           `envconfig:"uaa_origin"`                          // IGNITION_UAA_ORIGIN << REQUIRED
	ClientID          string           `envconfig:"api_client_id"`                       // IGNITION_API_CLIENT_ID << REQUIRED
	ClientSecret      string           `envconfig:"api_client_secret"`                   // IGNITION_API_CLIENT_SECRET << REQUIRED
	SkipTLSValidation bool             `envconfig:"skip_tls_validation" default:"false"` // IGNITION_SKIP_TLS_VALIDATION
	CC                cloudfoundry.API `ignored:"true"`                                  // Ignored
	UAA               uaa.API          `ignored:"true"`                                  // Ignored
}

// NewDeployment uses environment variables to populate a Deployment
func NewDeployment(name string) (*Deployment, error) {
	var d Deployment
	envconfig.Process(ignition, &d)
	if cfenv.IsRunningOnCF() {
		c, err := cfenv.Current()
		if err != nil {
			return nil, err
		}
		s, err := c.Services.WithName(name)
		if err != nil {
			return nil, err
		}

		systemDomain, ok := s.CredentialString("system_domain")
		if ok && strings.TrimSpace(systemDomain) != "" {
			d.SystemDomain = systemDomain
		}

		skipTLSValidation, ok := s.CredentialString("skip_tls_validation")
		if ok {
			if b, err := strconv.ParseBool(skipTLSValidation); err == nil {
				d.SkipTLSValidation = b
			}
		}

		uaaOrigin, ok := s.CredentialString("uaa_origin")
		if ok && strings.TrimSpace(uaaOrigin) != "" {
			d.UAAOrigin = uaaOrigin
		}

		clientID, ok := s.CredentialString("api_client_id")
		if ok && strings.TrimSpace(clientID) != "" {
			d.ClientID = clientID
		}
		clientSecret, ok := s.CredentialString("api_client_secret")
		if ok && strings.TrimSpace(clientSecret) != "" {
			d.ClientSecret = clientSecret
		}
	}
	if strings.TrimSpace(d.SystemDomain) == "" {
		return nil, errors.New("system_domain is required")
	}
	d.ParseSystemDomain()
	d.UAAOrigin = strings.TrimSpace(d.UAAOrigin)
	if d.UAAOrigin == "" {
		return nil, errors.New("uaa_origin is required")
	}
	if strings.TrimSpace(d.ClientID) == "" {
		return nil, errors.New("api_client_id is required")
	}
	if strings.TrimSpace(d.ClientSecret) == "" {
		return nil, errors.New("api_client_secret is required")
	}

	uaaConfig := d.Config()
	uaaClient := uaaConfig.Client(context.Background())

	uaaAPI := &uaa.Client{
		URL:          d.UAAURL,
		ClientID:     d.ClientID,
		ClientSecret: d.ClientSecret,
		Client:       uaaClient,
		Config:       uaaConfig,
	}

	config := &cfclient.Config{
		ApiAddress:        d.APIURL,
		ClientID:          d.ClientID,
		ClientSecret:      d.ClientSecret,
		UserAgent:         "ignition-api",
		SkipSslValidation: d.SkipTLSValidation,
		HttpClient:        uaaClient,
		TokenSource:       uaaConfig.TokenSource(context.Background()),
	}

	d.CC = &cfclient.Client{
		Config: *config,
		Endpoint: cfclient.Endpoint{
			TokenEndpoint: uaaConfig.TokenURL,
		},
	}

	d.UAA = uaaAPI
	return &d, nil
}

// Config builds an oauth2.Config for the Deployment
func (d *Deployment) Config() *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     d.ClientID,
		ClientSecret: d.ClientSecret,
		TokenURL:     fmt.Sprintf("%s/oauth/token", d.UAAURL),
		Scopes:       []string{"cloud_controller.admin", "scim.write", "scim.read"},
	}
}

// URL builds a URL from the system domain, prepending the given string, if
// supplied
func (d *Deployment) URL(s string) *url.URL {
	if !strings.Contains(d.SystemDomain, "://") {
		d.SystemDomain = fmt.Sprintf("https://%s", d.SystemDomain)
	}
	u, err := url.Parse(d.SystemDomain)
	if err != nil {
		log.Println(err)
		return nil
	}

	if strings.TrimSpace(s) != "" && u.Hostname() != "127.0.0.1" && strings.ToLower(u.Hostname()) != "localhost" {
		u.Host = fmt.Sprintf("%s.%s", s, u.Host)
	}
	return u
}

// ParseSystemDomain sets the AppsURL, APIURL and UAAURL using the system domain
func (d *Deployment) ParseSystemDomain() {
	d.SystemDomain = strings.TrimSpace(d.SystemDomain)
	a := d.URL("apps")
	if a != nil {
		d.AppsURL = a.String()
	}
	a = d.URL("api")
	if a != nil {
		d.APIURL = a.String()
	}
	a = d.URL("login")
	if a != nil {
		d.UAAURL = a.String()
	}
}
