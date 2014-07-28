package solr

import (
	"fmt"
	"net/url"
	"strings"
)


type Query struct {
	params url.Values
}

func NewQuery() *Query {
	q := new(Query)
	q.params = url.Values{}
	return q
}

func(q *Query) AddParam(k string, v string) {
	q.params.Add(k, v)
}

func(q *Query) RemoveParam(k string) {
	q.params.Del(k)
}

func (q *Query) String() string {
	return q.params.Encode()
}

type Search struct {
	queries []*Query
	conn    *Connection
	start   int
	rows    int
	debug   string
}

func NewSearch(c *Connection, q *Query) * Search {
	s := new(Search)
	if q != nil {
		s.AddQuery(q)
	}
	if c != nil {
		s.conn = c
	}
	return s
}


func (s *Search) Query() *Query {
	q := NewQuery()
	s.AddQuery(q)
	return q
}

func (s *Search) AddQuery(q *Query) {
	s.queries = append(s.queries, q)
}

func (s *Search) QueryString() string {

	query := []string{"wt=json"}
	
	if s.start > 0 {
		query = append(query, fmt.Sprintf("start=%d", s.start))
	}

	if s.rows > 0 {
		query = append(query, fmt.Sprintf("rows=%d", s.rows))
	}

	if s.debug != "" {
		query = append(query, fmt.Sprintf("debug=%s&indent=true", s.debug))
	}

	if len(s.queries) > 0 {
		for _, v := range s.queries {
			query = append(query, v.String())
		}
	}

	return strings.Join(query, "&")
}

func (s *Search) Result() (*SolrResult, error) {
	resp, err := s.conn.Select(s.QueryString())
	if err != nil {
		return nil, err
	}
	result := &SolrResult{response: resp.response}
	result.results = new(Collection)
	result.status = resp.status
	if resp.status == 0 {
		if response, ok := resp.response["response"].(map[string]interface{}); ok {
			result.results.numFound = int(response["numFound"].(float64))
			result.results.start = int(response["start"].(float64))
			if docs, ok := response["docs"].([]interface{}); ok {
				for _, v := range docs {
					result.results.docs = append(result.results.docs, Document(v.(map[string]interface {})))
				}
			}
		}
	}
	
	return result, nil
}
