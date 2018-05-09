package cloudfoundry

import (
	"fmt"
	"net/url"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

// ISOSegmentQuerier is used to query a Cloud Controller API for isolation segments
type ISOSegmentQuerier interface {
	ListIsolationSegmentsByQuery(query url.Values) ([]cfclient.IsolationSegment, error)
}

// ISOSegmentIDForName gets the isolation segment ID for the given iso segment name
func ISOSegmentIDForName(name string, iq ISOSegmentQuerier) (string, error) {
	q := url.Values{}
	q.Set("names", name)
	isoSegments, err := iq.ListIsolationSegmentsByQuery(q)
	if err != nil {
		return "", err
	}
	for _, is := range isoSegments {
		if strings.EqualFold(is.Name, name) {
			return is.GUID, nil
		}
	}
	return "", fmt.Errorf("Could not find isolation segment named [%s]", name)
}
