
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
	q.AddParam("testing", []string{"test"})
	s := si.Search(q)
	fmt.Println(s.QueryString())
}