package solr

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"encoding/base64"
	"strings"
)

func logPrintBasicAuth(req *http.Request) {
	authData, ok := req.Header["Authorization"]
	if ok == false {
		return
	}
	auth := strings.SplitN(authData[0], " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return
	}
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	log.Printf("Basic auth: %v", pair)
}

func logRequest(req *http.Request) {
	if os.Getenv("MOCK_LOGGING") != "" {
		log.Printf("RequestURI: %s", req.RequestURI)
		logPrintBasicAuth(req)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err.Error())
		}
		log.Println(string(body))
		log.Println(req.Header)
	}
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

// For commands that no need of specific response
func mockSuccessCommand(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":5}}`)
}

func mockSuccessGrouped(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
		  "responseHeader":{
		    "status":0,
		    "QTime":2,
		    "params":{
		      "q":"*:*",
		      "group.field":"id",
		      "group":"true",
		      "wt":"json"}},
		  "grouped":{
		    "id":{
		      "matches":2,
		      "groups":[{
		          "groupValue":"test_id_100",
		          "doclist":{"numFound":1,"start":0,"docs":[
		              {
		                "id":"test_id_100",
		                "title":["add sucess 100"],
		                "_version_":1475623982992457728}]
		          }},
		        {
		          "groupValue":"test_id_101",
		          "doclist":{"numFound":1,"start":0,"docs":[
		              {
		                "id":"test_id_101",
		                "title":["add sucess 101"],
		                "_version_":1475623982995603456}]
		          }}]}}}`)
}

func mockSuccessStrangeGrouped(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{
		  "responseHeader":{
		    "status":0,
		    "QTime":2,
		    "params":{
		      "q":"*:*",
		      "group.field":"id",
		      "group":"true",
		      "wt":"json"}}}`)
}

func mockSuccessXML(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `<response>
						<lst name="responseHeader">
						<int name="status">0</int>
						<int name="QTime">8</int>
						</lst>
						</response>`)
}

func mockCoreAdmin(w http.ResponseWriter, req *http.Request) {
	logRequest(req)
	io.WriteString(w, `{"responseHeader":{"status":0,"QTime":50}}`)
}

func mockStartServer() {
	http.HandleFunc("/success/core0/select/", mockSuccessSelect)
	http.HandleFunc("/fail/core0/select/", mockFailSelect)
	http.HandleFunc("/facet_counts/core0/select/", mockSuccessSelectFacet)
	http.HandleFunc("/highlight/core0/select/", mockSuccessSelectHighlight)

	http.HandleFunc("/standalonecommit/core0/update/", mockSuccessStandaloneCommit)
	http.HandleFunc("/add/core0/update/", mockSuccessAdd)
	http.HandleFunc("/delete/core0/update/", mockSuccessDelete)

	http.HandleFunc("/command/core0/update/", mockSuccessCommand)
	http.HandleFunc("/xml/core0/update/", mockSuccessXML)
	http.HandleFunc("/grouped/core0/select/", mockSuccessGrouped)
	http.HandleFunc("/noresponse/core0/select/", mockSuccessStrangeGrouped)
	http.HandleFunc("/solr/admin/cores", mockCoreAdmin)
	
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

