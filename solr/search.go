package solr

import (
	"fmt"
	"net/url"
)

type Query struct {
	params *url.Values
}

func NewQuery() *Query {
	q := new(Query)
	q.params = &url.Values{}
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
// Example: id:100
func (qq *Query) Q(q string) {
	qq.params.Add("q", q)
}

// sort parameter http://wiki.apache.org/solr/CommonQueryParameters
// Example: geodist() asc
func (q *Query) Sort(sort string) {
	q.params.Add("sort", sort)
}

// fq (Filter Query) http://wiki.apache.org/solr/CommonQueryParameters
// Example: popularity:[10 TO *]
func (q *Query) FilterQuery(fq string) {
	q.params.Add("fq", fq)
}

// fl (Field List ) parameter http://wiki.apache.org/solr/CommonQueryParameters
// Example: id,name,decsription
func (q *Query) FieldList(fl string) {
	q.params.Add("fl", fl)
}

// geofilt - The distance filter http://wiki.apache.org/solr/SpatialSearch
// Output example: fq={!geofilt pt=45.15,-93.85 sfield=store d=5}
func (q *Query) Geofilt(latitude float64, longitude float64, sfield string, distance float64) {
	q.params.Add("fq", fmt.Sprintf("{!geofilt pt=%#v,%#v sfield=%s d=%#v}", latitude, longitude, sfield, distance))
}

// defType http://wiki.apache.org/solr/CommonQueryParameters
// Example: dismax
func (q *Query) DefType(defType string) {
	q.params.Add("defType", defType)
}

// bf (Boost Functions) parameter http://wiki.apache.org/solr/DisMaxQParserPlugin
// Example: ord(popularity)^0.5 recip(rord(price),1,1000,1000)^0.3
func (q *Query) BoostFunctions(bf string) {
	q.params.Add("bf", bf)
}

// bq (Boost Query) parameter http://wiki.apache.org/solr/DisMaxQParserPlugin
func (q *Query) BoostQuery(bq string) {
	q.params.Add("bq", bq)
}

// qf (Query Fields) parameter http://wiki.apache.org/solr/DisMaxQParserPlugin
// Example: features^20.0+text^0.3
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

// Return query params including debug and indent if Debug is set
func (s *Search) QueryParams() *url.Values {

	if s.query == nil {
		s.query = NewQuery()
	}

	if s.Debug != "" {
		s.query.params.Set("debug", s.Debug)
		s.query.params.Set("indent", "true")
	}

	return s.query.params
}

// QueryString return a query string of all queries except wt=json
func (s *Search) QueryString() string {
	return s.QueryParams().Encode()
}

// Wrapper for connection.Resource which will add wt=json automatically
// One can use this to query to /solr/{CORE}/{RESOURCE} example /solr/collection1/select
// This can be useful when you use an search component that is not supported in this package
func (s *Search) Resource(resource string, params *url.Values) (*SolrResponse, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("No connection found for making request to solr")
	}
	return s.conn.Resource(resource, params)
}

// Result will create a StandardResultParser if no parser specified.
// parser must be an implementation of ResultParser interface
func (s *Search) Result(parser ResultParser) (*SolrResult, error) {
	resp, err := s.Resource("select", s.QueryParams())
	if err != nil {
		return nil, err
	}
	if parser == nil {
		parser = new(StandardResultParser)
	}
	return parser.Parse(resp)
}

// This method is for making query to MoreLikeThisHandler
// See http://wiki.apache.org/solr/MoreLikeThisHandler
func (s *Search) MoreLikeThis(parser MltResultParser) (*SolrMltResult, error) {
	resp, err := s.Resource("mlt", s.QueryParams())
	if err != nil {
		return nil, err
	}
	if parser == nil {
		parser = new(MoreLikeThisParser)
	}
	return parser.Parse(resp)
}
