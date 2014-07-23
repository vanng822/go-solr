package solr

import "testing"

import "fmt"

func TestSolrQueryAddParam(t *testing.T) {
	
	q := NewQuery()
	q.AddParam("qf", "some qf")
		
	if q.String() != "qf=some+qf" {
		t.Errorf("Expected to be: 'some qf'")
	}
	
	fmt.Println(q.String())
}

func TestSolrSearchMultipleQuery(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	s := NewSearch(nil, q)
	q2 := NewQuery()
	q2.AddParam("testing", "testing 2")
	s.AddQuery(q2)
	if s.QueryString() != "testing=test&testing=testing+2" {
		t.Errorf("Expected to be: 'testing=test&testing=testing+2'")
	}
}

func TestSolrQueryRemoveParam(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	q.AddParam("testing2", "testing 2")
	// random order in for loop of range on map
	if q.String() != "testing=test&testing2=testing+2" && q.String() != "testing2=testing+2&testing=test" {
		t.Errorf("Expected to be: 'testing=test&testing2=testing+2' or 'testing2=testing+2&testing=test'")
	}
	q.RemoveParam("testing2")
	if q.String() != "testing=test" {
		t.Errorf("Expected to be: 'testing=test'")
	}
}
