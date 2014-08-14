package solr

import (
	"fmt"
	"net/url"
	"testing"
)

func TestSolrDocument(t *testing.T) {
	d := Document{"id": "test_id", "title": "test title"}
	if d.Has("id") == false {
		t.Errorf("Has id expected to be true")
	}

	if d.Has("not_exist") == true {
		t.Errorf("Has not_exist expected to be false")
	}

	if d.Get("title").(string) != "test title" {
		t.Errorf("title expected to have value 'test title'")
	}

	d.Set("new_title", "new title")
	if d.Get("new_title").(string) != "new title" {
		t.Errorf("new_title expected to have value 'new title'")
	}
	
	if d.Get("not_exist") != nil {
		t.Errorf("Get not_exist key should return nil but got '%s'", d.Get("not_exist"))
	}
}

func TestSolrSuccessSelect(t *testing.T) {
	go mockStartServer()

	si, err := NewSolrInterface("http://127.0.0.1:12345/success", "core0")

	if err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}

	q := NewQuery()
	q.AddParam("q", "*:*")
	s := si.Search(q)

	res, err := s.Result(nil)

	if err != nil {
		t.Errorf("cannot seach %s", err)
	}

	if res.Status != 0 {
		t.Errorf("Status expected to be 0")
	}

	if res.Results.NumFound != 1 {
		t.Errorf("results.numFound expected to be 1")
	}

	if res.Results.Start != 0 {
		t.Errorf("results.start expected to be 0")
	}

	if len(res.Results.Docs) != 1 {
		t.Errorf("len of .docs should be 1")
	}

	if res.Results.Docs[0].Get("id").(string) != "change.me" {
		t.Errorf("id of first document should be change.me")
	}
}

func TestSolrFailSelect(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/fail", "core0")

	if err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}

	q := NewQuery()
	q.AddParam("q", "*:*")
	s := si.Search(q)

	parser := new(StandardResultParser)
	res, err := s.Result(parser)

	if err != nil {
		t.Errorf("cannot seach %s", err)
	}

	if res.Status != 400 {
		t.Errorf("Status expected to be 400")
	}
	expectedMsg := "no field name specified in query and no default specified via 'df' param"
	msg, ok := res.Error["msg"].(string)
	if ok != true {
		t.Errorf("error expected to have a message")
	}

	if msg != expectedMsg {
		t.Errorf("Error msg expected to be '%s' but got '%s'", expectedMsg, msg)
	}

	if res.Results.NumFound != 0 {
		t.Errorf("results.numFound expected to be 0")
	}

	if res.Results.Start != 0 {
		t.Errorf("results.start expected to be 0")
	}

	if len(res.Results.Docs) != 0 {
		t.Errorf("len of .docs should be 0")
	}

}

func TestSolrFacetSelect(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/facet_counts", "core0")

	if err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}

	q := NewQuery()
	q.AddParam("q", "*:*")
	q.AddParam("facet", "true")
	q.AddParam("facet.field", "id")

	s := si.Search(q)
	parser := new(StandardResultParser)
	res, err := s.Result(parser)

	if err != nil {
		t.Errorf("cannot seach %s", err)
	}

	if res.Status != 0 {
		t.Errorf("Status expected to be 0 but got %d", res.Status)
	}

	if res.Results.NumFound != 4 {
		t.Errorf("results.numFound expected to be 4 but got %d", res.Results.NumFound)
	}

	if res.Results.Start != 0 {
		t.Errorf("results.start expected to be 0 but got %d", res.Results.Start)
	}

	if len(res.Results.Docs) != 4 {
		t.Errorf("len of .docs should be 4 but got %d", len(res.Results.Docs))
	}

	third_doc := res.Results.Docs[2]

	if third_doc.Get("id") != "change.me3" {
		t.Errorf("id of third document expected to be 'change.me3' but got '%s'", third_doc.Get("id"))
	}

	if _, ok := res.FacetCounts["facet_fields"]; ok == false {
		t.Errorf("results.facet_counts.facet_fields expected")
		return
	}

	facet_fields := res.FacetCounts["facet_fields"].(map[string]interface{})
	id, ok := facet_fields["id"]

	if ok == false {
		t.Errorf("results.facet_counts.facet_fields.id expected")
		return
	}

	id_len := len(id.([]interface{}))

	if id_len != 6 {
		t.Errorf("results.facet_counts.facet_fields.id.len expected be 6 but got %d", id_len)
	}
	
	if _, ok := res.ResponseHeader["params"]; ok == false {
		t.Errorf("ResponseHeader should contain param key")
	}
}

func TestSolrHighlightSelect(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/highlight", "core0")

	if err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}

	q := NewQuery()
	q.AddParam("q", "*:*")
	q.AddParam("hl", "true")

	s := si.Search(q)
	parser := new(StandardResultParser)
	res, err := s.Result(parser)

	if err != nil {
		t.Errorf("cannot seach %s", err)
	}

	if res.Status != 0 {
		t.Errorf("Status expected to be 0 but got %d", res.Status)
	}

	if res.Results.NumFound != 4 {
		t.Errorf("results.numFound expected to be 4 but got %d", res.Results.NumFound)
	}

	if res.Results.Start != 0 {
		t.Errorf("results.start expected to be 0 but got %d", res.Results.Start)
	}

	if len(res.Results.Docs) != 4 {
		t.Errorf("len of .docs should be 4 but got %d", len(res.Results.Docs))
	}

	third_doc := res.Results.Docs[2]

	if third_doc.Get("id") != "change.me3" {
		t.Errorf("id of third document expected to be 'change.me3' but got '%s'", third_doc.Get("id"))
	}

	_, ok := res.Highlighting["change.me"]

	if ok == false {
		t.Errorf("results.facet_counts.facet_fields.id expected")
		return
	}
}

func TestSolrResultLoopSelect(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/facet_counts", "core0")
	if err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}
	q := NewQuery()
	q.AddParam("q", "*:*")
	q.AddParam("facet", "true")
	q.AddParam("facet.field", "id")
	s := si.Search(q)
	res, err := s.Result(nil)

	if err != nil {
		t.Errorf("Should not have an error here, skip assertions below. Please fix!")
		return
	}

	if cap(res.Results.Docs) != 4 {
		t.Errorf("Capacity expected to be 4 but got '%d'", cap(res.Results.Docs))
	}

	if len(res.Results.Docs) != 4 {
		t.Errorf("len of .docs should be 4 but got %d", len(res.Results.Docs))
	}

	for i, doc := range res.Results.Docs {
		if doc.Has("id") == false {
			t.Errorf("Document %d doesn't contain id", i)
		}
	}

	for i := 0; i < len(res.Results.Docs); i++ {
		if res.Results.Docs[i].Has("id") == false {
			t.Errorf("Document %d doesn't contain id", i)
		}
	}

}

func TestSolrSuccessStandaloneCommit(t *testing.T) {

	si, err := NewSolrInterface("http://127.0.0.1:12345/standalonecommit", "core0")

	if err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}

	res, err := si.Commit()

	if err != nil {
		t.Errorf("cannot commit %s", err)
	}

	if res.Success != true {
		t.Errorf("success expected to be true")
	}
}

func TestMakeAddChunks(t *testing.T) {
	docs := make([]Document, 0, 100)
	for i := 0; i < 500; i++ {
		docs = append(docs, Document{"id": fmt.Sprintf("test_id_%d", i), "title": fmt.Sprintf("add sucess %d", i)})
	}
	chunks := makeAddChunks(docs, 100)
	expected_len := 5
	if len(chunks) != expected_len {
		t.Errorf("Chunks length expected to be '%d' but got '%d'", expected_len, len(chunks))
	}

	d := chunks[0]["add"].([]Document)[0]

	if d.Get("id") != "test_id_0" {
		t.Errorf("First element in first chunk should have id test_id_0 ")
	}

	d = chunks[1]["add"].([]Document)[0]

	if d.Get("id") != "test_id_100" {
		t.Errorf("First element in second chunk should have id test_id_100 ")
	}

	chunks = makeAddChunks(docs, 50)
	expected_len = 10
	if len(chunks) != expected_len {
		t.Errorf("Chunks length expected to be '%d' but got '%d'", expected_len, len(chunks))
	}

	chunks = makeAddChunks(docs, 301)
	expected_len = 2
	if len(chunks) != expected_len {
		t.Errorf("Chunks length expected to be '%d' but got '%d'", expected_len, len(chunks))
	}

	d = chunks[1]["add"].([]Document)[0]

	if d.Get("id") != "test_id_301" {
		t.Errorf("First element in second chunk should have id test_id_301 ")
	}

	if len(chunks[1]["add"].([]Document)) != 199 {
		t.Errorf("Last chunk should have length of 199 documents")
	}
}
func TestAdd(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/add", "core0")
	if err != nil {
		t.Errorf(err.Error())
	}
	si.SetBasicAuth("test", "post")
	
	docs := make([]Document, 0, 5)
	for i := 0; i < 5; i++ {
		docs = append(docs, Document{"id": fmt.Sprintf("test_id_%d", i), "title": fmt.Sprintf("add sucess %d", i)})
	}
	res, _ := si.Add(docs, 0, nil)
	res2, _ := si.Commit()
	// not sure what we can test here but at least run and see thing flows
	if res == nil {
		t.Errorf("Add response should not be nil")
	}

	if res2 == nil {
		t.Errorf("Commit response should not be nil")
	}
}

func TestDelete(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/delete", "core0")
	if err != nil {
		t.Errorf(err.Error())
	}

	res, _ := si.Delete(map[string]interface{}{"query": "id:test_id_1 OR id:test_id_2", "commitWithin": "500"}, nil)

	// not sure what we can test here but at least run and see thing flows
	if res == nil {
		t.Errorf("Delete response should not be nil")
	}

	params := &url.Values{}
	params.Add("commitWithin", "500")

	res2, _ := si.Delete(map[string]interface{}{"query": "*:*"}, params)

	// not sure what we can test here but at least run and see thing flows
	if res2 == nil {
		t.Errorf("Delete response should not be nil")
	}
}

func TestXMLResponse(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/xml", "core0")
	if err != nil {
		t.Errorf(err.Error())
	}

	res, err := si.DeleteAll()

	if err == nil {
		t.Errorf("Error should be not nil since response is not json format")
	}

	if err.Error() != "invalid character '<' looking for beginning of value" {
		t.Errorf("Expected error message 'invalid character '<' looking for beginning of value' but got '%s'", err.Error())
	}

	if res != nil {
		t.Errorf("Response should be nil since response is not json format")
	}
}

func TestRollback(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/command", "core0")
	if err != nil {
		t.Errorf(err.Error())
	}

	res, _ := si.Rollback()

	// not sure what we can test here but at least run and see thing flows
	if res == nil {
		t.Errorf("Rollback response should not be nil")
	}
}

func TestOptimize(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/command", "core0")
	if err != nil {
		t.Errorf(err.Error())
	}

	res, _ := si.Optimize(nil)

	// not sure what we can test here but at least run and see thing flows
	if res == nil {
		t.Errorf("Optimize response should not be nil")
	}
	params := &url.Values{}
	params.Add("maxSegments", "10")
	params.Add("waitFlush", "false")
	res2, _ := si.Optimize(params)

	// not sure what we can test here but at least run and see thing flows
	if res2 == nil {
		t.Errorf("Optimize response should not be nil")
	}
}

func TestGrouped(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/grouped", "core0")
	if err != nil {
		t.Errorf(err.Error())
	}
	
	q := NewQuery()
	q.AddParam("q", "*:*")
	q.AddParam("group", "true")
	q.AddParam("group.field", "id")
	
	s := si.Search(q)
	si.SetBasicAuth("test", "get")
	res, err := s.Result(nil)
	
	if err != nil {
		t.Errorf("Error should be nil")
		return
	}
	if _, ok := res.Grouped["id"]; ok == false {
		t.Errorf("should have key id in grouped")
	}
	
	if res.Results.Docs != nil {
		t.Errorf("Docs should be nil")
	}
}

func TestNoResponseGrouped(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/noresponse", "core1")
	if err != nil {
		t.Errorf(err.Error())
	}
	
	q := NewQuery()
	q.AddParam("q", "*:*")
	q.AddParam("group", "true")
	q.AddParam("group.field", "id")
	si.SetCore("core0")
	s := si.Search(q)
	
	_, err = s.Result(nil)
	
	if err == nil {
		t.Errorf("Error should not be nil")
	}
	expected := `Standard parser can only parse solr response with response object,
					ie response.response and response.response.docs. Or grouped response
					Please use other parser or implement your own parser`
	if err.Error() != expected {
		t.Errorf("expected error '%s' but got '%s'", expected, err.Error())
	}
}

/*
func TestRealAdd(t *testing.T) {
	fmt.Println("test_real")
	si, err := NewSolrInterface("http://localhost:8983/solr", "collection1")
	if err != nil {
		t.Errorf(err.Error())
	}

	docs := make([]Document,0,100)
	for i := 0; i < 100; i++ {
		docs = append(docs, Document{"id": fmt.Sprintf("test_id_%d", i), "title": fmt.Sprintf("add sucess %d", i)})
	}
	res, _ := si.Add(docs, 0, nil)

	res2, _ := si.Commit()
	si.Delete(map[string]interface{}{"query":"*:*"}, nil)
	//si.DeleteAll()
	si.Rollback()
	//si.Optimize(nil)
	params := &url.Values{}
	params.Add("maxSegments", "10")
	params.Add("waitFlush", "false")
	si.Optimize(params)
	fmt.Println(res.Result)
	fmt.Println(res2.Result)

	s := si.Search(nil)
	query := NewQuery()
	query.AddParam("q", "title:add sucess 1")
	s.SetQuery(query)
	r, err := s.Result(nil)

	fmt.Println(r.Results)
}
*/

/*
func TestRealDelete(t *testing.T) {
	fmt.Println("test_real")
	si, err := NewSolrInterface("http://localhost:8983/solr", "collection1")
	if err != nil {
		t.Errorf(err.Error())
	}
	params := &url.Values{}
	params.Add("commitWithin", "500")
	res, _ := si.Delete(map[string]interface{}{ "query":"id:test_id_0 OR id:test_id_1"}, params)

	fmt.Println(res.Result)
}
*/
/*
func TestRealDeleteAll(t *testing.T) {
	fmt.Println("test_real")
	si, err := NewSolrInterface("http://localhost:8983/solr", "collection1")
	if err != nil {
		t.Errorf(err.Error())
	}

	res, _ := si.DeleteAll()

	fmt.Println(res.Result)
}
*/

func TestCoreAdminCoresAction(t *testing.T) {
	
	ca, _ := NewCoreAdmin("http://127.0.0.1:12345/solr")
	
	params := &url.Values{}
	params.Add("core", "core0")
	res, err := ca.Action("RELOAD", params)
	if err != nil {
		t.Errorf("Should not be an error")
	}
	if res.Status != 0 {
		t.Errorf("Status expected to be 0 but got '%d'", res.Status)
	}
}

func TestCoreAdminCoresActionWrappers(t *testing.T) {
	
	ca, _ := NewCoreAdmin("http://127.0.0.1:12345/solr")
	
	// Status
	res, err := ca.Status("")
	if err != nil {
		t.Errorf("Should not be an error")
	}
	if res.Status != 0 {
		t.Errorf("Status expected to be 0 but got '%d'", res.Status)
	}
	
	res, err = ca.Status("core0")
	if err != nil {
		t.Errorf("Should not be an error")
	}
	if res.Status != 0 {
		t.Errorf("Status expected to be 0 but got '%d'", res.Status)
	}
	// Swap
	
	res, err = ca.Swap("core0", "core1")
	if err != nil {
		t.Errorf("Should not be an error")
	}
	if res.Status != 0 {
		t.Errorf("Status expected to be 0 but got '%d'", res.Status)
	}
	
	// Reload
	res, err = ca.Reload("core0")
	if err != nil {
		t.Errorf("Should not be an error")
	}
	if res.Status != 0 {
		t.Errorf("Status expected to be 0 but got '%d'", res.Status)
	}
	
	// Unload
	res, err = ca.Unload("core0")
	if err != nil {
		t.Errorf("Should not be an error")
	}
	if res.Status != 0 {
		t.Errorf("Status expected to be 0 but got '%d'", res.Status)
	}
	
	// Rename
	res, err = ca.Rename("core0", "core5")
	
	if err != nil {
		t.Errorf("Should not be an error")
	}
	if res.Status != 0 {
		t.Errorf("Status expected to be 0 but got '%d'", res.Status)
	}
	// Split
	res, err = ca.Split("core0", "core1")
	
	if err == nil {
		t.Errorf("Should be an error")
	}
	
	expected := "You must specify at least 2 target cores"
	
	if err.Error() != expected {
		t.Errorf("expcted '%s' but got '%s'", expected, err.Error())
	}
	
	res, err = ca.Split("core0", "core1", "core2")
	
	if err != nil {
		t.Errorf("Should not be an error")
	}
	
	if res.Status != 0 {
		t.Errorf("Status expected to be 0 but got '%d'", res.Status)
	}
}

func TestSupportedAction(t *testing.T) {

	ca, _ := NewCoreAdmin("http://127.0.0.1:12345/solr")

	params := &url.Values{}
	actions := []string{ "CREATE", "mergeindexes"}
	for _, action := range actions {
		_, err := ca.Action(action, params)
		if err != nil {
			t.Errorf("Should not be an error but got '%s'", err.Error())
		}
	}
}