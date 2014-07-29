package solr

import (
	"io"
	"log"
	"net/http"
)

func mockSuccessSelect(w http.ResponseWriter, req *http.Request) {
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

func mockFailSelect(w http.ResponseWriter, req *http.Request) {
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

func mockStartServer() {
	http.HandleFunc("/success/select/", mockSuccessSelect)
	http.HandleFunc("/fail/select/", mockFailSelect)
	http.HandleFunc("/facet_counts/select/", mockSuccessSelectFacet)

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
