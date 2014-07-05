package solr

import (
	"net/url"
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
	u, err := url.Parse(solrUrl)
	if err != nil {
		return nil, err
	}
	conn := &Connection{url: u}

	si := &SolrInterface{conn: conn, format: "json"}

	return si, nil
}

func (self *SolrInterface) Search() {

}

func (self *SolrInterface) Add(docs []string) {

}

func (self *SolrInterface) Delete(data string) (*UpdateResponse, error) {
	// prepare delete message here
	message := data
	return self.conn.Update(message)
}

func (self *SolrInterface) Update(data string) (*UpdateResponse, error) {
	// prepare message
	message := data

	return self.conn.Update(message)
}

func (self *SolrInterface) Commit() (*UpdateResponse, error) {
	return self.conn.Commit()
}

func (self *SolrInterface) Optimize() (*UpdateResponse, error) {
	return self.conn.Optimize()
}

func (self *SolrInterface) Rollback() (*UpdateResponse, error) {
	return self.conn.Rollback()
}
