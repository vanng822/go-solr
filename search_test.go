package solr

import "testing"
import "fmt"

func TestSolrQuery(t *testing.T) {
	q := Query{
		params: QueryParams{
			"q": "{!type=dismax qf=myfield v='solr rocks'}",
		},
	}

	//q := NewQuery()

	q.params["bf"] = "something"

	q.AddParam("qf", "some qf")
	q.AddParam("sbrm", "should be removed")
	q.RemoveParam("sbrm")
	fmt.Println(q.String())
}

func TestSolrSearch(t *testing.T) {
	//s := new(Search)
	//s.conn = &Connection
}
