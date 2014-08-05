package solr

import (
	"fmt"
	"testing"
	"net/url"
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
}

func TestSolrSuccessSelect(t *testing.T) {
	go mockStartServer()

	si, err := NewSolrInterface("http://127.0.0.1:12345/success")

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

	parser := new(StandardResultParser)
	res, err := s.Result(parser)

	if err != nil {
		t.Errorf("cannot seach %s", err)
	}

	if res.status != 400 {
		t.Errorf("Status expected to be 400")
	}
	expectedMsg := "no field name specified in query and no default specified via 'df' param"
	msg, ok := res.error["msg"].(string)
	if ok != true {
		t.Errorf("error expected to have a message")
	}

	if msg != expectedMsg {
		t.Errorf("Error msg expected to be '%s' but got '%s'", expectedMsg, msg)
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

func TestSolrFacetSelect(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/facet_counts")

	if err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}

	q := NewQuery()
	q.AddParam("q", "*:*")
	q.AddParam("facet", "true")
	q.AddParam("facet.field", "id")

	s := si.Search(q)
	fmt.Println(s.QueryString())
	parser := new(StandardResultParser)
	res, err := s.Result(parser)

	if err != nil {
		t.Errorf("cannot seach %s", err)
	}

	if res.status != 0 {
		t.Errorf("Status expected to be 0 but got %d", res.status)
	}

	if res.results.numFound != 4 {
		t.Errorf("results.numFound expected to be 4 but got %d", res.results.numFound)
	}

	if res.results.start != 0 {
		t.Errorf("results.start expected to be 0 but got %d", res.results.start)
	}

	if len(res.results.docs) != 4 {
		t.Errorf("len of .docs should be 4 but got %d", len(res.results.docs))
	}

	third_doc := res.results.docs[2]

	if third_doc.Get("id") != "change.me3" {
		t.Errorf("id of third document expected to be 'change.me3' but got '%s'", third_doc.Get("id"))
	}

	if _, ok := res.facet_counts["facet_fields"]; ok == false {
		t.Errorf("results.facet_counts.facet_fields expected")
		return
	}

	facet_fields := res.facet_counts["facet_fields"].(map[string]interface{})
	id, ok := facet_fields["id"]

	if ok == false {
		t.Errorf("results.facet_counts.facet_fields.id expected")
		return
	}

	id_len := len(id.([]interface{}))

	if id_len != 6 {
		t.Errorf("results.facet_counts.facet_fields.id.len expected be 6 but got %d", id_len)
	}
}

func TestSolrHighlightSelect(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/highlight")

	if err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}

	q := NewQuery()
	q.AddParam("q", "*:*")
	q.AddParam("hl", "true")

	s := si.Search(q)
	fmt.Println(s.QueryString())
	parser := new(StandardResultParser)
	res, err := s.Result(parser)

	if err != nil {
		t.Errorf("cannot seach %s", err)
	}

	if res.status != 0 {
		t.Errorf("Status expected to be 0 but got %d", res.status)
	}

	if res.results.numFound != 4 {
		t.Errorf("results.numFound expected to be 4 but got %d", res.results.numFound)
	}

	if res.results.start != 0 {
		t.Errorf("results.start expected to be 0 but got %d", res.results.start)
	}

	if len(res.results.docs) != 4 {
		t.Errorf("len of .docs should be 4 but got %d", len(res.results.docs))
	}

	third_doc := res.results.docs[2]

	if third_doc.Get("id") != "change.me3" {
		t.Errorf("id of third document expected to be 'change.me3' but got '%s'", third_doc.Get("id"))
	}

	_, ok := res.highlighting["change.me"]

	if ok == false {
		t.Errorf("results.facet_counts.facet_fields.id expected")
		return
	}
}

func TestSolrResultLoopSelect(t *testing.T) {
	si, err := NewSolrInterface("http://127.0.0.1:12345/facet_counts")
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

	if cap(res.results.docs) != 4 {
		t.Errorf("Capacity expected to be 4 but got '%d'", cap(res.results.docs))
	}

	if len(res.results.docs) != 4 {
		t.Errorf("len of .docs should be 4 but got %d", len(res.results.docs))
	}

	for i, doc := range res.results.docs {
		if doc.Has("id") == false {
			t.Errorf("Document %d doesn't contain id", i)
		}
	}

	for i := 0; i < len(res.results.docs); i++ {
		if res.results.docs[i].Has("id") == false {
			t.Errorf("Document %d doesn't contain id", i)
		}
	}

}

func TestSolrSuccessStandaloneCommit(t *testing.T) {

	si, err := NewSolrInterface("http://127.0.0.1:12345/standalonecommit")

	if err != nil {
		t.Errorf("Can not instance a new solr interface, err: %s", err)
	}

	res, err := si.Commit()

	if err != nil {
		t.Errorf("cannot commit %s", err)
	}

	if res.success != true {
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
	fmt.Println("test_real")
	si, err := NewSolrInterface("http://127.0.0.1:12345/add")
	if err != nil {
		t.Errorf(err.Error())
	}

	docs := make([]Document,0,5)
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
	fmt.Println("test_real")
	si, err := NewSolrInterface("http://127.0.0.1:12345/delete")
	if err != nil {
		t.Errorf(err.Error())
	}

	res, _ := si.Delete(map[string]interface{}{ "query":"id:test_id_1 OR id:test_id_2", "commitWithin":"500" }, nil)
	
	// not sure what we can test here but at least run and see thing flows
	if res == nil {
		t.Errorf("Delete response should not be nil")
	}
	
	params := &url.Values{}
	params.Add("commitWithin", "500")
	
	res2, _ := si.Delete(map[string]interface{}{ "query":"*:*"}, params)
	
	// not sure what we can test here but at least run and see thing flows
	if res2 == nil {
		t.Errorf("Delete response should not be nil")
	}
}


func TestRollback(t *testing.T) {
	fmt.Println("test_real")
	si, err := NewSolrInterface("http://127.0.0.1:12345/command")
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
	fmt.Println("test_real")
	si, err := NewSolrInterface("http://127.0.0.1:12345/command")
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

/*
func TestRealAdd(t *testing.T) {
	fmt.Println("test_real")
	si, err := NewSolrInterface("http://localhost:8983/solr/collection1")
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
	fmt.Println(res.result)
	fmt.Println(res2.result)
}
*/
/*
func TestRealDelete(t *testing.T) {
	fmt.Println("test_real")
	si, err := NewSolrInterface("http://localhost:8983/solr/collection1")
	if err != nil {
		t.Errorf(err.Error())
	}
	params := &url.Values{}
	params.Add("commitWithin", "500")
	res, _ := si.Delete(map[string]interface{}{ "query":"id:test_id_0 OR id:test_id_1"}, params)
	
	fmt.Println(res.result)
}
*/
/*
func TestRealDeleteAll(t *testing.T) {
	fmt.Println("test_real")
	si, err := NewSolrInterface("http://localhost:8983/solr/collection1")
	if err != nil {
		t.Errorf(err.Error())
	}
	
	res, _ := si.DeleteAll()
	
	fmt.Println(res.result)
}
*/
