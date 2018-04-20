package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/dghubble/sessions"
	"github.com/kelseyhightower/envconfig"
)

// Server is an HTTP/S web server
type Server struct {
	ServiceName   string         `envconfig:"config_servicename" default:"ignition-config"` // IGNITION_CONFIG_SERVICENAME
	Scheme        string         `envconfig:"scheme" default:"http"`                        // IGNITION_SCHEME
	Domain        string         `envconfig:"domain" default:"localhost"`                   // IGNITION_DOMAIN
	Port          int            `envconfig:"port" default:"3000"`                          // IGNITION_PORT
	ServePort     int            `envconfig:"serve_port" default:"3000"`                    // IGNITION_SERVE_PORT
	WebRoot       string         `ignored:"true"`                                           // Not configurable
	SessionSecret string         `envconfig:"session_secret"`                               // IGNITION_SESSION_SECRET << REQUIRED
	CompanyName   string         `envconfig:"company_name" default:"Your Company"`          // IGNITION_COMPANY_NAME
	SessionStore  sessions.Store `ignored:"true"`                                           // Not configurable
}

// NewServer uses environment variables to populate a Server
func NewServer() (*Server, error) {
	var s Server
	err := envconfig.Process(ignition, &s)
	if err != nil {
		return nil, err
	}
	err = s.ConfigureServer(s.ServiceName)
	if err != nil {
		return nil, err
	}
	root, _ := os.Getwd()
	s.ConfigureWebRoot(root)
	if strings.TrimSpace(s.SessionSecret) == "" {
		return nil, errors.New("session_secret is required")
	}
	s.SessionStore = sessions.NewCookieStore([]byte(s.SessionSecret), nil)
	return &s, nil
}

// ConfigureWebRoot ensures the webroot is set to appropriate values for local
// development and for use on Cloud Foundry
func (s *Server) ConfigureWebRoot(root string) {
	if cfenv.IsRunningOnCF() {
		s.WebRoot = root
	} else {
		s.WebRoot = filepath.Join(root, "web", "dist")
	}
}

// ConfigureServer ensures that the server will function correctly on Cloud
// Foundry
func (s *Server) ConfigureServer(name string) error {
	if !cfenv.IsRunningOnCF() {
		return nil
	}
	env, err := cfenv.Current()
	if err != nil {
		return err
	}
	s.Scheme = "https"
	s.Port = 443
	s.ServePort = env.Port
	service, err := env.Services.WithName(name)
	if err == nil && service != nil {
		domain, ok := service.CredentialString("domain")
		if ok && strings.TrimSpace(domain) != "" {
			s.Domain = domain
		}

		sessionSecret, ok := service.CredentialString("session_secret")
		if ok && strings.TrimSpace(sessionSecret) != "" {
			s.SessionSecret = sessionSecret
		}

		companyName, ok := service.CredentialString("company_name")
		if ok && strings.TrimSpace(companyName) != "" {
			s.CompanyName = companyName
		}
	}
	d := strings.TrimSpace(strings.ToLower(s.Domain))
	if d != "localhost" && d != "" {
		return nil
	}
	if len(env.ApplicationURIs) == 0 {
		return errors.New("ignition requires a route to function; please map a route")
	}
	s.Domain = env.ApplicationURIs[0]
	return nil
}
