
package solr

import (
	"fmt"
	"testing"
)

func TestSolr(t *testing.T) {
	si, err := NewSolrInterface("https://www.test.tld")
	
	if err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}
	q := NewQuery()
	q.AddParam("testing", "test")
	s := si.Search(q)
	q2 := NewQuery()
	q2.AddParam("testing", "testing 2")
	s.AddQuery(q2)
	
	fmt.Println(s.QueryString())
}