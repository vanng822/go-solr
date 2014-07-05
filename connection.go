package solr

import (
	"net/url"
)

type SelectResponse struct {
	results *Collection
	status  int
	qtime   int
}

type UpdateResponse struct {
	success bool
}

type ErrorResponse struct {
	message string
	status  int
}

type Connection struct {
	url *url.URL
}

func (self *Connection) Select(select_url string) (*SelectResponse, error) {
	return nil, nil
}

func (self *Connection) Update(data string) (*UpdateResponse, error) {
	return nil, nil
}

func (self *Connection) Commit() (*UpdateResponse, error) {
	return nil, nil
}

func (self *Connection) Optimize() (*UpdateResponse, error) {
	return nil, nil
}

func (self *Connection) Rollback() (*UpdateResponse, error) {
	return nil, nil
}
