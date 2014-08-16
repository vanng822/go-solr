package solr

import (
	"fmt"
	"net/url"
	"strings"
)

type Schema struct {
	url      *url.URL
	core     string
	username string
	password string
}

type SchemaResponse struct {
	Status int
	Result map[string]interface{}
}

// NewSchema will parse solrUrl and return a schema object, solrUrl must be a absolute url or path
func NewSchema(solrUrl, core string) (*Schema, error) {
	u, err := url.ParseRequestURI(solrUrl)
	if err != nil {
		return nil, err
	}

	return &Schema{url: u, core: core}, nil
}

// Set to a new core
func (s *Schema) SetCore(core string) {
	s.core = core
}

func (s *Schema) SetBasicAuth(username, password string) {
	s.username = username
	s.password = password
}

// See Get requests in https://wiki.apache.org/solr/SchemaRESTAPI for detail
func (s *Schema) Get(path string, params *url.Values) (*SchemaResponse, error) {
	var (
		r   []byte
		err error
	)
	if params == nil {
		params = &url.Values{}
	}

	params.Set("wt", "json")

	if s.core != "" {
		r, err = HTTPGet(fmt.Sprintf("%s/%s/schema/%s?%s", s.url.String(), s.core, strings.TrimLeft(path, "/"), params.Encode()), nil, s.username, s.password)
	} else {
		r, err = HTTPGet(fmt.Sprintf("%s/schema/%s?%s", s.url.String(), strings.TrimLeft(path, "/"), params.Encode()), nil, s.username, s.password)
	}
	if err != nil {
		return nil, err
	}
	resp, err := bytes2json(&r)
	if err != nil {
		return nil, err
	}

	return &SchemaResponse{Result: resp, Status: int(resp["responseHeader"].(map[string]interface{})["status"].(float64))}, nil
}
