package solr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func HTTPPost(path string, data *[]byte, headers [][]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", path, bytes.NewReader(*data))
	if len(headers) > 0 {
		for i := range headers {
			req.Header.Add(headers[i][0], headers[i][1])
		}
	}

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

func HTTPGet(url string, headers [][]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if len(headers) > 0 {
		for i := range headers {
			req.Header.Add(headers[i][0], headers[i][1])
		}
	}

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

func json2bytes(data map[string]interface{}) (*[]byte, error) {
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

type SelectResponse struct {
	/**
	responseHeader map[string]interface{}
	response       map[string]interface{}
	facet_counts   map[string]interface{}
	highlighting   map[string]interface{}
	grouped        map[string]interface{}
	debug          map[string]interface{}
	error          map[string]interface{}
	*/
	response map[string]interface{}
	// status quick access to status
	status int
}

type UpdateResponse struct {
	success bool
	result  map[string]interface{}
}


type Connection struct {
	url *url.URL
}

func NewConnection(solrUrl string) (*Connection, error) {
	u, err := url.Parse(solrUrl)
	if err != nil {
		return nil, err
	}

	return &Connection{url: u}, nil
}

func (c *Connection) Select(selectQuery string) (*SelectResponse, error) {
	r, err := HTTPGet(fmt.Sprintf("%s/select/?%s", c.url.String(), selectQuery), nil)
	if err != nil {
		return nil, err
	}
	resp, err := bytes2json(&r)
	if err != nil {
		return nil, err
	}

	result := SelectResponse{response: resp}
	result.status = int(resp["responseHeader"].(map[string]interface{})["status"].(float64))
	return &result, nil
}

func (c *Connection) Update(data map[string]interface{}) (*UpdateResponse, error) {
	b, err := json2bytes(data)
	if err != nil {
		return nil, err
	}
	r, err := HTTPPost(fmt.Sprintf("%s/update/", c.url.String()), b, nil)
	if err != nil {
		return nil, err
	}
	resp, err := bytes2json(&r)
	if err != nil {
		return nil, err
	}
	// check error in resp
	if hasError(resp) {
		return &UpdateResponse{success: false, result: resp}, nil
	}

	return &UpdateResponse{success: true, result: resp}, nil
}

func (c *Connection) Commit() (*UpdateResponse, error) {
	return nil, nil
}

func (c *Connection) Optimize() (*UpdateResponse, error) {
	return nil, nil
}

func (c *Connection) Rollback() (*UpdateResponse, error) {
	return nil, nil
}
