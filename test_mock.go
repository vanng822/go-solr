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

func mockFailSelect(w http.ResponseWriter, req *http.Request) {
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

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
