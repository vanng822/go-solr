package solr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var userAgent = fmt.Sprintf("Go-solr/%s (+https://github.com/vanng822/go-solr)", VERSION)

// HTTPPost make a POST request to path which also includes domain, headers are optional
func HTTPPost(path string, data *[]byte, headers [][]string, username, password string) ([]byte, error) {
	var (
		req *http.Request
		err error
	)

	client := &http.Client{}
	if data == nil {
		req, err = http.NewRequest("POST", path, nil)
	} else {
		req, err = http.NewRequest("POST", path, bytes.NewReader(*data))
	}

	if err != nil {
		return nil, err
	}

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	if len(headers) > 0 {
		for i := range headers {
			req.Header.Add(headers[i][0], headers[i][1])
		}
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// HTTPGet make a GET request to url, headers are optional
func HTTPGet(url string, headers [][]string, username, password string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	if len(headers) > 0 {
		for i := range headers {
			req.Header.Add(headers[i][0], headers[i][1])
		}
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	return body, nil
}

func bytes2json(data *[]byte) (map[string]interface{}, error) {
	var container interface{}

	err := json.Unmarshal(*data, &container)

	if err != nil {
		return nil, err
	}

	return container.(map[string]interface{}), nil
}

func json2bytes(data interface{}) (*[]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func hasError(response map[string]interface{}) bool {
	_, ok := response["error"]
	return ok
}

type Connection struct {
	url      *url.URL
	core     string
	username string
	password string
}

// NewConnection will parse solrUrl and return a connection object, solrUrl must be a absolute url or path
func NewConnection(solrUrl, core string) (*Connection, error) {
	u, err := url.ParseRequestURI(solrUrl)
	if err != nil {
		return nil, err
	}

	return &Connection{url: u, core: core}, nil
}

// Set to a new core
func (c *Connection) SetCore(core string) {
	c.core = core
}

func (c *Connection) SetBasicAuth(username, password string) {
	c.username = username
	c.password = password
}

func (c *Connection) Resource(source, queryString string) (*SolrResponse, error) {
	r, err := HTTPGet(fmt.Sprintf("%s/%s/%s?%s", c.url.String(), c.core, source, queryString), nil, c.username, c.password)
	if err != nil {
		return nil, err
	}
	resp, err := bytes2json(&r)
	if err != nil {
		return nil, err
	}

	result := SolrResponse{Response: resp}
	result.Status = int(resp["responseHeader"].(map[string]interface{})["status"].(float64))
	return &result, nil
}

// Update take optional params which can use to specify addition parameters such as commit=true
func (c *Connection) Update(data map[string]interface{}, params *url.Values) (*SolrUpdateResponse, error) {

	b, err := json2bytes(data)

	if err != nil {
		return nil, err
	}

	if params == nil {
		params = &url.Values{}
	}

	params.Set("wt", "json")

	r, err := HTTPPost(fmt.Sprintf("%s/%s/update/?%s", c.url.String(), c.core, params.Encode()), b, [][]string{{"Content-Type", "application/json"}}, c.username, c.password)

	if err != nil {
		return nil, err
	}
	resp, err := bytes2json(&r)
	if err != nil {
		return nil, err
	}
	// check error in resp
	if hasError(resp) {
		return &SolrUpdateResponse{Success: false, Result: resp}, nil
	}

	return &SolrUpdateResponse{Success: true, Result: resp}, nil
}
