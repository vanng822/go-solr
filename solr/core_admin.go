package solr

import (
	"fmt"
	"net/url"
)

type CoreAdmin struct {
	url      *url.URL
	username string
	password string
}

type CoreAdminResponse struct {
	// Quick access to responseHeader.status
	Status int
	// Raw response from solr core admin
	Result map[string]interface{}
}

// solrUrl should look like this http://0.0.0.0:8983/solr[/admin/cores] ie /admin/cores will append automatically
// when calling Action
func NewCoreAdmin(solrUrl string) (*CoreAdmin, error) {
	u, err := url.ParseRequestURI(solrUrl)
	if err != nil {
		return nil, err
	}

	return &CoreAdmin{url: u}, nil
}

// Set basic auth in case solr require login
func (ca *CoreAdmin) SetBasicAuth(username, password string) {
	ca.username = username
	ca.password = password
}

// Call to admin/cores endpoint, additional to params neccessary for this action can specified in params.
// No check is done for those params so check https://wiki.apache.org/solr/CoreAdmin for detail
// action is case sensitive
func (ca *CoreAdmin) Action(action string, params *url.Values) (*CoreAdminResponse, error) {
	switch action {
	case "STATUS":
	case "RELOAD":
	case "CREATE":
	case "RENAME":
	case "SWAP":
	case "UNLOAD":
	case "SPLIT":
	case "mergeindexes":
		params.Set("action", action)
	default:
		return nil, fmt.Errorf("Action '%s' not supported", action)
	}
	
	params.Set("wt", "json")
	
	r, err := HTTPGet(fmt.Sprintf("%s/admin/cores?%s", ca.url.String(), params.Encode()), nil, ca.username, ca.password)
	if err != nil {
		return nil, err
	}
	resp, err := bytes2json(&r)
	if err != nil {
		return nil, err
	}
	result := &CoreAdminResponse{Result: resp}
	result.Status = int(resp["responseHeader"].(map[string]interface{})["status"].(float64))
	return result, nil
}
