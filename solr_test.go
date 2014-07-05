
package solr

import "testing"


func SolrTest(t *testing.T) {

	if _, err := NewSolrInterface(""); err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}
}