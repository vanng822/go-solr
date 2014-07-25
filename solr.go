package solr

import (
)

type Document struct {
	fields map[string]interface{}
}

type Collection struct {
	docs     []Document
	start    int
	numFound int
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

func (si *SolrInterface) Add(docs []*Document) {
	
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
