package solr


import (
	"net/url"
)

type Document struct {
	fields map[string]interface{}
}

type Collection struct {
	docs []Document
	start int
	numFound int
}


type SolrInterface struct{
	format string
	conn *Connection
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

func (si *SolrInterface) Search() {

}

func (si *SolrInterface) Add(docs []string) {

}

func (si *SolrInterface) Delete(data string) (*UpdateResponse, error) {
	// prepare delete message here
	message := data
	return si.conn.Update(message)
}

func (si *SolrInterface) Update(data string) (*UpdateResponse, error) {
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

