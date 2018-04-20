package config

import (
	"strings"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pivotalservices/ignition/cloudfoundry"
	"github.com/pkg/errors"
)

const defaultQuota string = "default"

// Experimenter is the metadata required to vend a Cloud Foundry organization
// and space for developer experimentation
type Experimenter struct {
	OrgPrefix string `envconfig:"org_prefix" default:"ignition"`   // IGNITION_ORG_PREFIX
	SpaceName string `envconfig:"space_name" default:"playground"` // IGNITION_SPACE_NAME
	QuotaName string `envconfig:"quota_name" default:"ignition"`   // IGNITION_QUOTA_NAME
	QuotaID   string `ignored:"true"`
}

// NewExperimenter uses environment variables to populate an Experimenter
func NewExperimenter(name string, q cloudfoundry.QuotaQuerier) (*Experimenter, error) {
	var e Experimenter
	envconfig.Process(ignition, &e)
	if cfenv.IsRunningOnCF() {
		env, err := cfenv.Current()
		if err != nil {
			return nil, err
		}
		service, err := env.Services.WithName(name)
		if err == nil && service != nil {
			orgPrefix, ok := service.CredentialString("org_prefix")
			if ok && strings.TrimSpace(orgPrefix) != "" {
				e.OrgPrefix = orgPrefix
			}
			quotaName, ok := service.CredentialString("quota_name")
			if ok && strings.TrimSpace(quotaName) != "" {
				e.QuotaName = quotaName
			}
			spaceName, ok := service.CredentialString("space_name")
			if ok && strings.TrimSpace(spaceName) != "" {
				e.SpaceName = spaceName
			}
		}
	}
	e.OrgPrefix = strings.TrimSpace(e.OrgPrefix)
	e.QuotaName = strings.TrimSpace(e.QuotaName)
	e.SpaceName = strings.TrimSpace(e.SpaceName)

	if e.QuotaName == "" {
		e.QuotaName = defaultQuota
	}
	id, err := cloudfoundry.QuotaIDForName(e.QuotaName, q)
	if err != nil {
		var defaultErr error
		id, defaultErr = cloudfoundry.QuotaIDForName(defaultQuota, q)
		if defaultErr != nil {
			return nil, errors.Wrapf(err, "could not find quota id for quota with name [%s], nor for the default quota", e.QuotaName)
		}
		e.QuotaName = defaultQuota
	}
	e.QuotaID = id
	return &e, nil
}
