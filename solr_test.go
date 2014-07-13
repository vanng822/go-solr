
package solr

import "testing"


func TestSolr(t *testing.T) {

	if _, err := NewSolrInterface("https://www.test.tld"); err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}
}