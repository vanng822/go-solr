package solr

// ResultParser is interface for parsing result from response.
// The idea here is that application have possibility to parse.
// Or defined own parser with internal data structure to suite
// application's need
type ResultParser interface {
	Parse(response *SelectResponse) (*SolrResult, error)
}

type StandardResultParser struct {
}

func (parser *StandardResultParser) Parse(response *SelectResponse) (*SolrResult, error) {
	sr := &SolrResult{}
	sr.results = new(Collection)
	sr.status = response.status

	if response.status == 0 {
		parser.ParseResponse(response, sr)
		parser.ParseFacetCounts(response, sr)
		parser.ParseHighlighting(response, sr)
	} else {
		parser.ParseError(response, sr)
	}

	return sr, nil
}

func (parser *StandardResultParser) ParseError(response *SelectResponse, sr *SolrResult) {
	if error, ok := response.response["error"]; ok {
		sr.error = error.(map[string]interface{})
	}
}

// ParseResponse will assign result and build sr.docs if there is a response.
// If there is no response property in response it will panic
func (parser *StandardResultParser) ParseResponse(response *SelectResponse, sr *SolrResult) {
	if resp, ok := response.response["response"].(map[string]interface{}); ok {
		sr.results.numFound = int(resp["numFound"].(float64))
		sr.results.start = int(resp["start"].(float64))
		if docs, ok := resp["docs"].([]interface{}); ok {
			sr.results.docs = make([]Document, len(docs))
			for i, v := range docs {
				sr.results.docs[i] = Document(v.(map[string]interface{}))
			}
		}
	} else {
		panic(`Standard parser can only parse solr response with response object,
					ie response.response and response.response.docs.
					Please use other parser or implement your own parser`)
	}
}

// ParseFacetCounts will assign facet_counts to sr if there is one.
// No modification done here
func (parser *StandardResultParser) ParseFacetCounts(response *SelectResponse, sr *SolrResult) {
	if facetCounts, ok := response.response["facet_counts"]; ok {
		sr.facet_counts = facetCounts.(map[string]interface{})
	}
}

// ParseHighlighting will assign highlighting to sr if there is one.
// No modification done here
func (parser *StandardResultParser) ParseHighlighting(response *SelectResponse, sr *SolrResult) {
	if highlighting, ok := response.response["highlighting"]; ok {
		sr.highlighting = highlighting.(map[string]interface{})
	}
}
