package solr

import "testing"

func TestParseMoreLikeThisMatch(t *testing.T) {
	data := []byte(`{
			  "responseHeader":{
			    "status":0,
			    "QTime":4},
			  "match":{"numFound":200,"start":0,"docs":[
			      {
			        "id":"test_id_0",
			        "title":["add sucess 0"],
			        "_version_":1476345720316362752}]
			  },
			  "response":{"numFound":199,"start":0,"docs":[
			      {
			        "id":"test_id_1",
			        "title":["add sucess 1"],
			        "_version_":1476345720644567040},
			      {
			        "id":"test_id_2",
			        "title":["add sucess 2"],
			        "_version_":1476345720645615616},
			      {
			        "id":"test_id_3",
			        "title":["add sucess 3"],
			        "_version_":1476345720645615617}]
			  }}`)
			  
	resp, _ := bytes2json(&data)
	response := &SolrResponse{Response: resp, Status: 0}
	
	parser := MoreLikeThisParser{}
	
	res, _ := parser.Parse(response)
	
	if res.Match.Start != 0 {
		t.Errorf("res.Match.Start expected to be '0' but got '%d'", res.Match.Start)
	}
	
	if len(res.Match.Docs) != 1 {
		t.Errorf("res.Match.Docs should have '1' doc but got '%d'", len(res.Match.Docs))
	}
	expected := "test_id_0"
	if res.Match.Docs[0].Get("id") != expected {
		t.Errorf("title expected to be '%s' but got '%s'", expected, res.Match.Docs[0].Get("id"))
	}
	
	if res.Results.Start != 0 {
		t.Errorf("res.Match.Start expected to be '0' but got '%d'", res.Results.Start)
	}
	
	if len(res.Results.Docs) != 3 {
		t.Errorf("res.Results.Docs should have '3' doc but got '%d'", len(res.Results.Docs))
	}
	expected = "test_id_1"
	if res.Results.Docs[0].Get("id") != expected {
		t.Errorf("title expected to be '%s' but got '%s'", expected, res.Results.Docs[0].Get("id"))
	}
}


