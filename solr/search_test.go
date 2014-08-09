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
	expected := "testing=test&testing=testing+2"
	if res != expected {
		t.Errorf("Expected to be: '%s' but got '%s'", expected, res)
	}
}

func TestSolrSearchMultipleValueSearchQuery(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	s := NewSearch(nil, q)
	q.AddParam("testing", "testing 2")
	res := s.QueryString()
	expected := "testing=test&testing=testing+2&wt=json"
	if res != expected {
		t.Errorf("Expected to be: '%s' but got '%s'", expected, res)
	}
}

func TestSolrSearchSetQuery(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	s := NewSearch(nil, q)
	expected := "testing=test&wt=json"
	res := s.QueryString()
	if res != expected {
		t.Errorf("Expected to be: '%s' but got '%s'", expected, res)
	}
	q2 := NewQuery()
	q2.AddParam("testing", "test2")
	s.SetQuery(q2)
	
	expected = "testing=test2&wt=json"
	res = s.QueryString()
	
	if res != expected {
		t.Errorf("Expected to be: '%s' but got '%s'", expected, res)
	}
}

func TestSolrSearchDebugQuery(t *testing.T) {
	q := NewQuery()
	q.AddParam("testing", "test")
	s := NewSearch(nil, q)
	s.Debug = "true"
	res := s.QueryString()
	expected := "debug=true&indent=true&testing=test&wt=json"
	if res != expected {
		t.Errorf("Expected to be: '%s' but got '%s'", expected, res)
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
	if res != "testing=test&testing2=testing+2" {
		t.Errorf("Expected to be: 'testing=test&testing2=testing+2' but got %s", res)
	}
	q.RemoveParam("testing2")
	if q.String() != "testing=test" {
		t.Errorf("Expected to be: 'testing=test'")
	}
}

func TestQueryQ(t *testing.T) {
	q := NewQuery()
	q.Q("id:100")
	expected := "q=id%3A100"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQuerySort(t *testing.T) {
	q := NewQuery()
	q.Sort("geodist() desc")
	expected := "sort=geodist%28%29+desc"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryFilterQuery(t *testing.T) {
	q := NewQuery()
	q.FilterQuery("popularity:[10 TO *]")
	expected := "fq=popularity%3A%5B10+TO+%2A%5D"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryFieldList(t *testing.T) {
	q := NewQuery()
	q.FieldList("id,name,decsription")
	expected := "fl=id%2Cname%2Cdecsription"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryGeofilt(t *testing.T) {
	q := NewQuery()
	q.Geofilt(45.15, -93.85, "store", 5)
	expected := "fq=%7B%21geofilt+pt%3D45.15%2C-93.85+sfield%3Dstore+d%3D5%7D"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryDefType(t *testing.T) {
	q := NewQuery()
	q.DefType("func")
	expected := "defType=func"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryBoostFunctions(t *testing.T) {
	q := NewQuery()
	q.BoostFunctions("recip(rord(myfield),1,2,3)")
	expected := "bf=recip%28rord%28myfield%29%2C1%2C2%2C3%29"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryBoostQuery(t *testing.T) {
	q := NewQuery()
	q.BoostQuery("cat:electronics^5.0")
	expected := "bq=cat%3Aelectronics%5E5.0"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}

func TestQueryQueryField(t *testing.T) {
	q := NewQuery()
	q.QueryFields("features^20.0+text^0.3")
	expected := "qf=features%5E20.0%2Btext%5E0.3"
	result := q.String()
	if result != expected {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}
