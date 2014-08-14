package solr

import (
	"fmt"
	"net/url"
	"strings"
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

// Call to admin/cores endpoint, additional params neccessary for this action can specified in params.
// No check is done for those params so check https://wiki.apache.org/solr/CoreAdmin for detail
func (ca *CoreAdmin) Action(action string, params *url.Values) (*CoreAdminResponse, error) {
	switch strings.ToUpper(action) {
	case "STATUS":
		params.Set("action", "STATUS")
	case "RELOAD":
		params.Set("action", "RELOAD")
	case "CREATE":
		params.Set("action", "CREATE")
	case "RENAME":
		params.Set("action", "RENAME")
	case "SWAP":
		params.Set("action", "SWAP")
	case "UNLOAD":
		params.Set("action", "UNLOAD")
	case "SPLIT":
		params.Set("action", "SPLIT")
	case "MERGEINDEXES":
		params.Set("action", "mergeindexes")
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

// pass empty string as core if you want status of all cores.
// See https://wiki.apache.org/solr/CoreAdmin#STATUS
func (ca *CoreAdmin) Status(core string) (*CoreAdminResponse, error) {
	params := &url.Values{}

	if core != "" {
		params.Add("core", core)
	}

	return ca.Action("STATUS", params)
}

// Swap one core with other core.
// See https://wiki.apache.org/solr/CoreAdmin#SWAP
func (ca *CoreAdmin) Swap(core, other string) (*CoreAdminResponse, error) {
	params := &url.Values{}
	params.Add("core", core)
	params.Add("other", other)
	return ca.Action("SWAP", params)
}

// Reload a core, see https://wiki.apache.org/solr/CoreAdmin#RELOAD
func (ca *CoreAdmin) Reload(core string) (*CoreAdminResponse, error) {
	params := &url.Values{}
	params.Add("core", core)
	return ca.Action("RELOAD", params)
}

// Unload a core, see https://wiki.apache.org/solr/CoreAdmin#UNLOAD
// If you want to use those flag deleteIndex, deleteDataDir, deleteInstanceDir
// Please use Action-method with those params specified, like ca.Action("UNLOAD", params)
func (ca *CoreAdmin) Unload(core string) (*CoreAdminResponse, error) {
	params := &url.Values{}
	params.Add("core", core)
	return ca.Action("UNLOAD", params)
}
