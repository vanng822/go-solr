package solr

import "testing"
import "fmt"

func TestSolrSearch(t *testing.T) {
	q := Query{
		params: QueryParams{
			"facet.field": []string{"accepts_4x4s", "accepts_bicycles"},
			"facet":       []string{"true"},
			"q":		[]string{"{!type=dismax qf=myfield v='solr rocks'}"},
			
		},
	}
	q.params["bf"] = []string{"something"}

	fmt.Println(q.String())
}
