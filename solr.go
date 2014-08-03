package solr

import (
	"fmt"
	"math"
	"net/url"
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

func makeAddChunks(docs []Document, chunk int) []map[string]interface{} {
	if chunk < 1 {
		chunk = 100
	}
	docs_len := len(docs)
	num_chunk := int(math.Ceil(float64(docs_len) / float64(chunk)))
	doc_counter := 0
	result := make([]map[string]interface{}, num_chunk)
	for i := 0; i < num_chunk; i++ {
		add := make([]Document, 0, chunk)
		for j := 0; j < chunk; j++ {
			if doc_counter >= docs_len {
				break
			}
			add = append(add, docs[doc_counter])
			doc_counter++
		}
		result[i] = map[string]interface{}{"add": add}
	}
	return result
}

func (si *SolrInterface) Add(docs []Document, chunk int, params *url.Values) (*UpdateResponse, error) {
	result := &UpdateResponse{success: true}
	responses := map[string]interface{}{}
	chunks := makeAddChunks(docs, chunk)

	for i := 0; i < len(chunks); i++ {
		res, err := si.Update(chunks[i], params)
		if err != nil {
			return nil, err
		}
		result.success = result.success && res.success
		responses[fmt.Sprintf("chunk_%d", i+1)] = map[string]interface{}{
			"result":  res.result,
			"success": res.success,
			"total":   len(chunks[i]["add"].([]Document))}
	}
	result.result = responses
	return result, nil
}

// Delete take data of type map and optional params which can use to specify addition parameters such as commit=true
func (si *SolrInterface) Delete(data map[string]interface{}, params *url.Values) (*UpdateResponse, error) {
	// prepare delete message here
	message := data
	return si.conn.Update(message, params)
}

// Update take data of type map and optional params which can use to specify addition parameters such as commit=true
func (si *SolrInterface) Update(data map[string]interface{}, params *url.Values) (*UpdateResponse, error) {
	// prepare message
	message := data

	return si.conn.Update(message, params)
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
