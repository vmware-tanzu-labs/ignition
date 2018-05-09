package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pivotalservices/ignition/cloudfoundry"
	"github.com/pkg/errors"
)

const defaultQuota string = "default"
const defaultIsolationSegment string = "shared"

// Experimenter is the metadata required to vend a Cloud Foundry organization
// and space for developer experimentation
type Experimenter struct {
	OrgPrefix              string        `envconfig:"org_prefix" default:"ignition"`          // IGNITION_ORG_PREFIX
	OrgCountUpdateInterval time.Duration `envconfig:"org_count_update_interval" default:"1m"` // IGNITION_ORG_COUNT_UPDATE_INTERVAL
	SpaceName              string        `envconfig:"space_name" default:"playground"`        // IGNITION_SPACE_NAME
	QuotaName              string        `envconfig:"quota_name" default:"ignition"`          // IGNITION_QUOTA_NAME
	QuotaID                string        `ignored:"true"`
	ISOSegmentName         string        `envconfig:"iso_segment_name" default:"shared"` // IGNITION_ISO_SEGMENT_NAME
	ISOSegmentID           string        `ignored:"true"`
}

// NewExperimenter uses environment variables to populate an Experimenter
func NewExperimenter(name string, qq cloudfoundry.QuotaQuerier, iq cloudfoundry.ISOSegmentQuerier) (*Experimenter, error) {
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
			updateInterval, ok := service.CredentialString("org_count_update_interval")
			if ok && strings.TrimSpace(updateInterval) != "" {
				d, err := time.ParseDuration(updateInterval)
				if err != nil {
					log.Println(fmt.Sprintf("[WARN] [%s] is an invalid time.Duration, defaulting org update interval to 1m", updateInterval))
				} else {
					e.OrgCountUpdateInterval = d
				}
			}
			quotaName, ok := service.CredentialString("quota_name")
			if ok && strings.TrimSpace(quotaName) != "" {
				e.QuotaName = quotaName
			}
			spaceName, ok := service.CredentialString("space_name")
			if ok && strings.TrimSpace(spaceName) != "" {
				e.SpaceName = spaceName
			}
			isoSegmentName, ok := service.CredentialString("iso_segment_name")
			if ok && strings.TrimSpace(isoSegmentName) != "" {
				e.ISOSegmentName = isoSegmentName
			}
		}
	}
	e.OrgPrefix = strings.TrimSpace(e.OrgPrefix)
	e.QuotaName = strings.TrimSpace(e.QuotaName)
	e.SpaceName = strings.TrimSpace(e.SpaceName)
	e.ISOSegmentName = strings.TrimSpace(e.ISOSegmentName)

	if e.QuotaName == "" {
		e.QuotaName = defaultQuota
	}
	quotaID, err := cloudfoundry.QuotaIDForName(e.QuotaName, qq)
	if err != nil {
		var defaultErr error
		quotaID, defaultErr = cloudfoundry.QuotaIDForName(defaultQuota, qq)
		if defaultErr != nil {
			return nil, errors.Wrapf(err, "could not find quota id for quota with name [%s], nor for the default quota", e.QuotaName)
		}
		e.QuotaName = defaultQuota
	}
	e.QuotaID = quotaID

	if e.ISOSegmentName == "" {
		e.ISOSegmentName = defaultIsolationSegment
	}
	isoSegmentID, err := cloudfoundry.ISOSegmentIDForName(e.ISOSegmentName, iq)
	if err != nil {
		return nil, err
	}
	e.ISOSegmentID = isoSegmentID
	return &e, nil
}
