package config

const ignition string = "ignition"

// Ignition is the configuration required for Ignition to function
type Ignition struct {
	Server       *Server
	Deployment   *Deployment
	Experimenter *Experimenter
	Authorizer   *Authorizer
}

// New builds configuration for Ignition using the environment and an associated
// Cloud Foundry user provided service named `ignition-config`, optionally
// making use of a bound p-identity service instance named `ignition`
func New() (*Ignition, error) {
	i := &Ignition{}
	s, err := NewServer()
	if err != nil {
		return nil, err
	}
	i.Server = s
	a, err := NewAuthorizer(s.ServiceName)
	if err != nil {
		return nil, err
	}
	i.Authorizer = a
	d, err := NewDeployment(s.ServiceName)
	if err != nil {
		return nil, err
	}
	i.Deployment = d
	e, err := NewExperimenter(s.ServiceName, i.Deployment.CC)
	if err != nil {
		return nil, err
	}
	i.Experimenter = e
	return i, nil
}
