package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pivotalservices/ignition/cloudfoundry"
	"github.com/pkg/errors"
)

// OrganizationHandler retrieves or creates the user's development organization
func OrganizationHandler(appsURL string, orgPrefix string, quotaID string, spaceName string, a cloudfoundry.API) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID, accountName, err := userInfoFromContext(req.Context())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if strings.TrimSpace(userID) == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		orgName := OrganizationName(orgPrefix, accountName)
		org, err := FindOrgForUser(orgName, appsURL, userID, quotaID, a)
		if err != nil {
			switch err.(type) {
			case OrgNotFoundError:
				org, err = CreateOrgForUser(orgName, appsURL, userID, quotaID, spaceName, a)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusNotFound)
					return
				}
			default:
				log.Println(err)
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(org)
	}
	return http.HandlerFunc(fn)
}

// OrgNotFoundError indicates that an org cannot be found for the user
type OrgNotFoundError string

func (o OrgNotFoundError) Error() string {
	return fmt.Sprintf("organization %s not found", o)
}

// CreateOrgForUser creates an org, a default space, and creates or retreieves
// the user and then assigns that user to org manager, org auditor, space manager,
// space developer, and space auditor roles
func CreateOrgForUser(name string, appsURL string, userID string, quotaID string, spaceName string, a cloudfoundry.API) (*cloudfoundry.Organization, error) {
	// create the user if needed
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("cannot create an org without a valid userID")
	}

	// check for org existence by name

	// create the org
	org, err := cloudfoundry.CreateOrg(name, appsURL, quotaID, a)
	if err != nil {
		return nil, err
	}

	// assign the user to org roles
	_, err = a.AssociateOrgUser(org.GUID, userID)
	if err != nil {
		log.Println(err)
	}
	_, err = a.AssociateOrgManager(org.GUID, userID)
	if err != nil {
		log.Println(err)
	}
	_, err = a.AssociateOrgAuditor(org.GUID, userID)
	if err != nil {
		log.Println(err)
	}

	// create the space and assign the user to all space roles
	err = cloudfoundry.CreateSpace(spaceName, org.GUID, userID, a)
	if err != nil {
		log.Println(err)
	}

	// return the org
	return org, nil
}

// FindOrgForUser returns an orgNotFoundError if the org is not found, and a
// single org with a name or quota match, when it exists
func FindOrgForUser(name string, appsURL string, userID string, quotaID string, a cloudfoundry.OrganizationQuerier) (*cloudfoundry.Organization, error) {
	o, err := cloudfoundry.OrgsForUserID(userID, appsURL, a)
	if err != nil {
		return nil, errors.Wrapf(err, "could not find orgs for user id: [%s]", userID)
	}

	if len(o) == 0 {
		return nil, OrgNotFoundError(name)
	}

	var quotaMatches []cloudfoundry.Organization
	for i := range o {
		if strings.EqualFold(quotaID, o[i].QuotaDefinitionGUID) {
			quotaMatches = append(quotaMatches, o[i])
		}
		if strings.EqualFold(name, o[i].Name) {
			return &o[i], nil
		}
	}

	if len(quotaMatches) == 0 {
		return nil, OrgNotFoundError(name)
	}

	return &quotaMatches[0], nil
}

// Name returns the organization name for the user's development organization
func OrganizationName(orgPrefix string, accountName string) string {
	orgPrefix = strings.ToLower(orgPrefix)
	accountName = strings.ToLower(accountName)
	if strings.Contains(accountName, "@") {
		components := strings.Split(accountName, "@")
		return fmt.Sprintf("%s-%s", orgPrefix, components[0])
	}

	if strings.Contains(accountName, "\\") {
		components := strings.Split(accountName, "\\")
		return fmt.Sprintf("%s-%s", orgPrefix, components[1])
	}

	return fmt.Sprintf("%s-%s", orgPrefix, accountName)
}
