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
	conn *Connection
}

func NewSolrInterface(solrUrl string) (*SolrInterface, error) {
	c, err := NewConnection(solrUrl)
	if err != nil {
		return nil, err
	}
	return &SolrInterface{conn: c}, nil
}

func (si *SolrInterface) Search(q *Query) *Search {
	s := NewSearch(si.conn, q)

	return s
}

// makeAddChunks splits the documents into chunks. If chunk_size is less than one it will be default to 100
func makeAddChunks(docs []Document, chunk_size int) []map[string]interface{} {
	if chunk_size < 1 {
		chunk_size = 100
	}
	docs_len := len(docs)
	num_chunk := int(math.Ceil(float64(docs_len) / float64(chunk_size)))
	doc_counter := 0
	chunks := make([]map[string]interface{}, num_chunk)
	for i := 0; i < num_chunk; i++ {
		add := make([]Document, 0, chunk_size)
		for j := 0; j < chunk_size; j++ {
			if doc_counter >= docs_len {
				break
			}
			add = append(add, docs[doc_counter])
			doc_counter++
		}
		chunks[i] = map[string]interface{}{"add": add}
	}
	return chunks
}

// Add will insert documents in batch of chunk_size. success is false as long as one chunk failed. 
// The result in UpdateResponse is summery of response from all chunks
// with key chunk_%d
func (si *SolrInterface) Add(docs []Document, chunk_size int, params *url.Values) (*UpdateResponse, error) {
	result := &UpdateResponse{success: true}
	responses := map[string]interface{}{}
	chunks := makeAddChunks(docs, chunk_size)

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
// Only one delete statement is supported, ie data can be { "id":"ID" }
// If you want to delete more docs use { "query":"QUERY" }
// Extra params can specify in params or in data such as { "query":"QUERY", "commitWithin":"500" }
func (si *SolrInterface) Delete(data map[string]interface{}, params *url.Values) (*UpdateResponse, error) {
	message := map[string]interface{}{"delete": data}
	return si.Update(message, params)
}

// DeleteAll will remove all documents and commit
func (si *SolrInterface) DeleteAll() (*UpdateResponse, error) {
	params := &url.Values{}
	params.Add("commit", "true")
	return si.Delete(map[string]interface{}{"query": "*:*"}, params)
}

// Update take data of type map and optional params which can use to specify addition parameters such as commit=true
func (si *SolrInterface) Update(data map[string]interface{}, params *url.Values) (*UpdateResponse, error) {
	if si.conn == nil {
		return nil, fmt.Errorf("No connection found for making request to solr")
	}
	return si.conn.Update(data, params)
}

func (si *SolrInterface) Commit() (*UpdateResponse, error) {
	params := &url.Values{}
	params.Add("commit", "true")
	return si.Update(map[string]interface{}{}, params)
}

func (si *SolrInterface) Optimize(params *url.Values) (*UpdateResponse, error) {
	if params == nil {
		params = &url.Values{}
	}
	params.Set("optimize", "true")
	return si.Update(map[string]interface{}{}, params)
}

// Rollback rollbacks all add/deletes made to the index since the last commit
// This should use with caution
// See https://wiki.apache.org/solr/UpdateXmlMessages#A.22rollback.22
func (si *SolrInterface) Rollback() (*UpdateResponse, error) {
	return si.Update(map[string]interface{}{"rollback": map[string]interface{}{}}, nil)
}
