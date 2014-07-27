package solr

import (
	"fmt"
	"testing"
)

func TestSolrSuccessSelect(t *testing.T) {
	go mockStartServer()

	si, err := NewSolrInterface("http://127.0.0.1:12345/success")

	if err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}
	
	q := NewQuery()
	q.AddParam("q", "*:*")
	s := si.Search(q)
	
	res, err := s.Result()

	if err != nil {
		t.Errorf("cannot seach %s", err)
	}
	
	if res.status != 0 {
		t.Errorf("Status expected to be 0")
	}
	
	if res.results.numFound != 1 {
		t.Errorf("results.numFound expected to be 1")
	}
	
	if res.results.start != 0 {
		t.Errorf("results.start expected to be 0")
	}
	
	if len(res.results.docs) != 1 {
		t.Errorf("len of .docs should be 1")
	}
	
	if res.results.docs[0].Get("id").(string) != "change.me" {
		t.Errorf("id of first document should be change.me")
	}
	
	fmt.Println(" ")
}


func TestSolrFailSelect(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/fail")

	if err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}
	
	q := NewQuery()
	q.AddParam("q", "*:*")
	s := si.Search(q)
	
	res, err := s.Result()

	if err != nil {
		t.Errorf("cannot seach %s", err)
	}
	
	if res.status != 400 {
		t.Errorf("Status expected to be 400")
	}
	
	if res.results.numFound != 0 {
		t.Errorf("results.numFound expected to be 0")
	}
	
	if res.results.start != 0 {
		t.Errorf("results.start expected to be 0")
	}
	
	if len(res.results.docs) != 0 {
		t.Errorf("len of .docs should be 0")
	}
	
	fmt.Println(" ")
}
