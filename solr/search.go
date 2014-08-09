package solr

import (
	"fmt"
	"net/url"
)

type Query struct {
	params url.Values
}

func NewQuery() *Query {
	q := new(Query)
	q.params = url.Values{}
	return q
}

func (q *Query) AddParam(k string, v string) {
	q.params.Add(k, v)
}

func (q *Query) RemoveParam(k string) {
	q.params.Del(k)
}

func (q *Query) GetParam(k string) string {
	return q.params.Get(k)
}

func (q *Query) SetParam(k string, v string) {
	q.params.Set(k, v)
}

// q parameter http://wiki.apache.org/solr/CommonQueryParameters
func (qq *Query) Q(q string) {
	qq.params.Add("q", q)
}

// sort parameter http://wiki.apache.org/solr/CommonQueryParameters
// geodist() asc
func (q *Query) Sort(sort string) {
	q.params.Add("sort", sort)
}

// fq (Filter Query) http://wiki.apache.org/solr/CommonQueryParameters
// popularity:[10 TO *]
func (q *Query) FilterQuery(fq string) {
	q.params.Add("fq", fq)
}

// fl (Field List ) parameter http://wiki.apache.org/solr/CommonQueryParameters
// id,name,decsription
func (q *Query) FieldList(fl string) {
	q.params.Add("fl", fl)
}

// geofilt - The distance filter http://wiki.apache.org/solr/SpatialSearch
// output example: fq={!geofilt pt=45.15,-93.85 sfield=store d=5}
func (q *Query) Geofilt(latitude float64, longitude float64, sfield string, distance float64) {
	q.params.Add("fq", fmt.Sprintf("{!geofilt pt=%#v,%#v sfield=%s d=%#v}", latitude, longitude, sfield, distance))
}

// defType http://wiki.apache.org/solr/CommonQueryParameters
func (q *Query) DefType(defType string) {
	q.params.Add("defType", defType)
}

// bf (Boost Functions) parameter http://wiki.apache.org/solr/DisMaxQParserPlugin
func (q *Query) BoostFunctions(bf string) {
	q.params.Add("bf", bf)
}

// bq (Boost Query) parameter http://wiki.apache.org/solr/DisMaxQParserPlugin
func (q *Query) BoostQuery(bq string) {
	q.params.Add("bq", bq)
}

// bf (Query Fields) parameter http://wiki.apache.org/solr/DisMaxQParserPlugin
func (q *Query) QueryFields(qf string) {
	q.params.Add("qf", qf)
}

func (q *Query) Start(start int) {
	q.params.Set("start", fmt.Sprintf("%d", start))
}

func (q *Query) Rows(rows int) {
	q.params.Set("rows", fmt.Sprintf("%d", rows))
}

func (q *Query) String() string {
	return q.params.Encode()
}

type Search struct {
	query *Query
	conn  *Connection
	Debug string
}

// NewSearch takes c and q as optional
func NewSearch(c *Connection, q *Query) *Search {
	s := new(Search)
	if q != nil {
		s.SetQuery(q)
	}
	
	if c != nil {
		s.conn = c
	}
	return s
}

// SetQuery will replace old query with new query q
func (s *Search) SetQuery(q *Query) {
	s.query = q
}

// QueryString return a query string of all queries, including start, rows, debug and wt=json.
// wt is always json
func (s *Search) QueryString() string {
	
	if s.query == nil {
		s.query = NewQuery()
	}
	
	s.query.params.Set("wt", "json")

	if s.Debug != "" {
		s.query.params.Set("debug", s.Debug)
		s.query.params.Set("indent", "true")
	}
	
	return s.query.String()
}

// Result will create a StandardResultParser if no parser specified.
// parser must be an implementation of ResultParser interface
func (s *Search) Result(parser ResultParser) (*SolrResult, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("No connection found for making request to solr")
	}
	resp, err := s.conn.Select(s.QueryString())
	if err != nil {
		return nil, err
	}
	if parser == nil {
		parser = new(StandardResultParser)
	}
	return parser.Parse(resp)
}
