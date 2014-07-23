package solr

import (
	"fmt"
	"net/url"
	"strings"
)

type QueryParams map[string][]string

type Query struct {
	params QueryParams
}

func NewQuery() *Query {
	q := new(Query)
	q.params = QueryParams{}
	return q
}

func(q *Query) AddParam(k string, v []string) {
	q.params[k] = v
}

func(q *Query) RemoveParam(k string) {
	delete(q.params, k)
}

func (q *Query) String() string {
	query := url.Values{}

	if len(q.params) > 0 {
		for k, v := range q.params {
			l := len(v)
			for x := 0; x < l; x++ {
				query.Add(k, v[x])
			}
		}
	}

	return query.Encode()
}

type Search struct {
	queries []*Query
	conn    *Connection
	start   int
	rows    int
	debug   string
}

func NewSearch(q *Query) * Search {
	s := new(Search)
	if q != nil {
		s.AddQuery(q)
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

	query := []string{}

	if s.start > 0 {
		query = append(query, fmt.Sprintf("start=%d", s.start))
	}

	if s.rows > 0 {
		query = append(query, fmt.Sprintf("rows=%d", s.rows))
	}

	if s.debug == "on" {
		query = append(query, "debug=on")
	}

	if len(s.queries) > 0 {
		for _, v := range s.queries {
			query = append(query, v.String())
		}
	}

	return strings.Join(query, "&")
}

func (s *Search) Result() (*Collection, error) {
	res, err := s.conn.Select(s.QueryString())
	if err != nil {
		return nil, err
	}
	// TODO fetch collection
	_ = res
	
	return nil, nil
}
