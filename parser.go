package solr

// ResultParser is interface for parsing result from response
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
	}

	return sr, nil
}


func (parser *StandardResultParser) ParseResponse(response *SelectResponse, sr *SolrResult) {
	if response, ok := response.response["response"].(map[string]interface{}); ok {
		sr.results.numFound = int(response["numFound"].(float64))
		sr.results.start = int(response["start"].(float64))
		if docs, ok := response["docs"].([]interface{}); ok {
			for _, v := range docs {
				sr.results.docs = append(sr.results.docs, Document(v.(map[string]interface{})))
			}
		}
	} else {
		panic(`Standard parser can only parse solr response with response object,
					ie response.response and response.response.docs.
					Please use other parser or implement your own parser`)
	}
}

func (parser *StandardResultParser) ParseFacetCounts(response *SelectResponse, sr *SolrResult) {
	if facetCounts, ok := response.response["facet_counts"]; ok {
		sr.facet_counts = facetCounts.(map[string]interface{})
	}
}

func (parser *StandardResultParser) ParseHighlighting(response *SelectResponse, sr *SolrResult) {
	if highlighting, ok := response.response["highlighting"]; ok {
		sr.highlighting = highlighting.(map[string]interface{})
	}
}
