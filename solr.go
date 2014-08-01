package solr

import (
	"fmt"
)

type Document map[string]interface{}

// Has check if a key exist in document
func (d Document) Has(k string) bool {
	_, ok := d[k]
	return ok
}

// Get returns value of a key if key exists else panic
func (d Document) Get(k string) interface{} {
	if v, ok := d[k]; ok {
		return v
	}
	panic(fmt.Sprintf("Try to access field '%s' which does not exist", k))
}

// Set add a key/value to document
func (d Document) Set(k string, v interface{}) {
	d[k] = v
}

type Collection struct {
	docs     []Document
	start    int
	numFound int
}

type SolrResult struct {
	status       int         // status quick access to status
	results      *Collection // results parsed documents, basically response object
	facet_counts map[string]interface{}
	highlighting map[string]interface{}
	error        map[string]interface{}

	// grouped for grouping result, not supported for now
	grouped map[string]interface{}
}

type SolrInterface struct {
	format string
	conn   *Connection
}

func NewSolrInterface(solrUrl string) (*SolrInterface, error) {
	c, err := NewConnection(solrUrl)
	if err != nil {
		return nil, err
	}
	return &SolrInterface{conn: c, format: "json"}, nil
}

func (si *SolrInterface) Search(q *Query) *Search {
	s := NewSearch(si.conn, q)

	return s
}

func (si *SolrInterface) Add(docs []Document) (*UpdateResponse, error) {
	return nil, nil
}

func (si *SolrInterface) Delete(data map[string]interface{}) (*UpdateResponse, error) {
	// prepare delete message here
	message := data
	return si.conn.Update(message)
}

func (si *SolrInterface) Update(data map[string]interface{}) (*UpdateResponse, error) {
	// prepare message
	message := data

	return si.conn.Update(message)
}

func (si *SolrInterface) Commit() (*UpdateResponse, error) {
	return si.conn.Commit()
}

func (si *SolrInterface) Optimize() (*UpdateResponse, error) {
	return si.conn.Optimize()
}

func (si *SolrInterface) Rollback() (*UpdateResponse, error) {
	return si.conn.Rollback()
}
