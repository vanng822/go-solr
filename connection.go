
package solr

import (
	"net/url"
)

type SelectResponse struct {
	results *Collection
	status int
	qtime int
}

type UpdateResponse struct {
	success bool 
}

type ErrorResponse struct {
	message string
	status int
}

type Connection struct {
	url *url.URL
}

func (c *Connection) Select(select_url string) (*SelectResponse, error) {
	return nil, nil
}

func (c *Connection) Update(data string) (*UpdateResponse, error) {
	return nil, nil
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