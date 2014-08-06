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

func TestSolrSearchMultipleValueQuery(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	q.AddParam("testing", "testing 2")
	res := q.String()
	if res != "testing=test&testing=testing+2" {
		t.Errorf("Expected to be: 'testing=test&testing=testing+2' but got '%s'", res)
	}
}

func TestSolrSearchMultipleValueSearchQuery(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	s := NewSearch(nil, q)
	q.AddParam("testing", "testing 2")
	if s.QueryString() != "wt=json&testing=test&testing=testing+2" {
		t.Errorf("Expected to be: 'wt=json&testing=test&testing=testing+2'")
	}
}

func TestSolrSearchSetQuery(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	s := NewSearch(nil, q)
	if s.QueryString() != "wt=json&testing=test" {
		t.Errorf("Expected to be: 'wt=json&testing=test'")
	}
	q2 := NewQuery()
	q2.AddParam("testing", "test2")
	s.SetQuery(q2)
	
	if s.QueryString() != "wt=json&testing=test2" {
		t.Errorf("Expected to be: 'wt=json&testing=test2'")
	}
}

func TestSolrSearchDebugQuery(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	s := NewSearch(nil, q)
	s.Debug = "true"
	res := s.QueryString()
	if res != "wt=json&debug=true&indent=true&testing=test" {
		t.Errorf("Expected to be: 'wt=json&debug=true&indent=true&testing=test' but got '%s'", res)
	}
}

func TestSolrSearchWithoutConnection(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	s := NewSearch(nil, q)
	
	resp, err := s.Result(&StandardResultParser{})
	if resp != nil {
		t.Errorf("resp expected to be nil due to no connection is set")
	}
	if err == nil {
		t.Errorf("err expected to be not empty due to no connection is set")
	}
	expectedErrorMessage := "No connection found for making request to solr"
	
	if err.Error() != expectedErrorMessage {
		t.Errorf("The error message expecte to be '%s' but got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestSolrQueryRemoveParam(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	q.AddParam("testing2", "testing 2")
	// random order in for loop of range on map
	res := q.String()
	if res != "testing2=testing+2&testing=test" && res != "testing=test&testing2=testing+2" {
		t.Errorf("Expected to be: 'testing2=testing+2&testing=test' or 'testing=test&testing2=testing+2' but got %s", res)
	}
	q.RemoveParam("testing2")
	if q.String() != "testing=test" {
		t.Errorf("Expected to be: 'testing=test'")
	}
}
