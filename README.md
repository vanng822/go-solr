go-solr
=======


[![Build Status](https://travis-ci.org/vanng822/go-solr.svg?branch=master)](https://travis-ci.org/vanng822/go-solr)
[![GoDoc](https://godoc.org/github.com/vanng822/go-solr?status.svg)](https://godoc.org/github.com/vanng822/go-solr)

Solr v4

Json only

No schema checking

Features: Search, Add, Update, Delete, Commit, Rollback, Optimize


## Install

go get github.com/vanng822/go-solr

## Usage

    package main
    import "github.com/vanng822/go-solr"
    import "fmt"
  
    func main() {
      si, _ := solr.NewSolrInterface("http://localhost:8983/solr/collection1")
      query := solr.NewQuery()
      query.AddParam("q", "*:*")
      s := si.Search(query)
      r, _ := s.Result(nil)
      fmt.Println(r.Results.Docs)
    }
