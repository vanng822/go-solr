package solr

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func logRequest(req *http.Request) {
	log.Printf("RequestURI: %s", req.RequestURI)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}
	log.Println(string(body))
}

func mockSuccessSelect(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
		  "responseHeader":{
		    "status":0,
		    "QTime":1,
		    "params":{
		      "indent":"true",
		      "q":"*:*",
		      "wt":"json"}},
		  "response":{"numFound":1,"start":0,"docs":[
		      {
		        "id":"change.me",
		        "title":["change.me"],
		        "_version_":1474699756018073600}]
		  }}`)
}

func mockSuccessSelectFacet(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
				  "responseHeader":{
				    "status":0,
				    "QTime":10,
				    "params":{
				      "facet":"true",
				      "indent":"true",
				      "q":"*:*",
				      "facet.field":"id",
				      "wt":"json"}},
				  "response":{"numFound":4,"start":0,"docs":[
				      {
				        "id":"change.me",
				        "title":["change.me"],
				        "_version_":1474893319511212032},
				      {
				        "id":"change.me2",
				        "title":["change.me2"],
				        "_version_":1474893328448225280},
				      {
				        "id":"change.me3",
				        "title":["change.me3"],
				        "_version_":1474893336208736256},
				      {
				        "id":"change.me2",
				        "title":["change.me22"],
				        "_version_":1474893362047746048}]
				  },
				  "facet_counts":{
				    "facet_queries":{},
				    "facet_fields":{
				      "id":[
				        "change.me2",2,
				        "change.me",1,
				        "change.me3",1]},
				    "facet_dates":{},
				    "facet_ranges":{}}}`)
}

func mockSuccessSelectHighlight(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
					  "responseHeader":{
					    "status":0,
					    "QTime":0,
					    "params":{
					      "indent":"true",
					      "q":"*:*",
					      "wt":"json",
					      "hl":"true"}},
					  "response":{"numFound":4,"start":0,"docs":[
					      {
					        "id":"change.me",
					        "title":["change.me"],
					        "_version_":1474893319511212032},
					      {
					        "id":"change.me2",
					        "title":["change.me2"],
					        "_version_":1474893328448225280},
					      {
					        "id":"change.me3",
					        "title":["change.me3"],
					        "_version_":1474893336208736256},
					      {
					        "id":"change.me2",
					        "title":["change.me22"],
					        "_version_":1474893362047746048}]
					  },
					  "highlighting":{
					    "change.me":{},
					    "change.me2":{},
					    "change.me3":{},
					    "change.me2":{}}}`)
}

func mockFailSelect(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	w.WriteHeader(400)
	io.WriteString(w, `{
		  "responseHeader":{
		    "status":400,
		    "QTime":3,
		    "params":{
		      "indent":"true",
		      "q":"**",
		      "wt":"json"}},
		  "error":{
		    "msg":"no field name specified in query and no default specified via 'df' param",
		    "code":400}}`)
}

func mockSuccessStandaloneCommit(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":5}}`)
}

func mockSuccessAdd(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":5}}`)
}

func mockSuccessDelete(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":5}}`)
}

func mockStartServer() {
	http.HandleFunc("/success/select/", mockSuccessSelect)
	http.HandleFunc("/fail/select/", mockFailSelect)
	http.HandleFunc("/facet_counts/select/", mockSuccessSelectFacet)
	http.HandleFunc("/highlight/select/", mockSuccessSelectHighlight)

	http.HandleFunc("/standalonecommit/update/", mockSuccessStandaloneCommit)
	http.HandleFunc("/add/update/", mockSuccessAdd)
	http.HandleFunc("/delete/update/", mockSuccessDelete)
	
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
