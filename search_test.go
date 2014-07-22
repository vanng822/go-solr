package solr

import "testing"
import "fmt"

func TestSolrQuery(t *testing.T) {
	q := Query{
		params: QueryParams{
			"facet.field": []string{"accepts_4x4s", "accepts_bicycles"},
			"facet":       []string{"true"},
			"q":		[]string{"{!type=dismax qf=myfield v='solr rocks'}"},
			
		},
	}
	q.params["bf"] = []string{"something"}
	
	q.AddParam("qf", []string{"some qf"})
	q.AddParam("sbrm", []string{"should be removed"})
	q.RemoveParam("sbrm")
	fmt.Println(q.String())
}

func TestSolrSearch(t *testing.T) {
	//s := new(Search)
	//s.conn = &Connection
}


