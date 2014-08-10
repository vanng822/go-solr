package solr
import (
	"fmt"
)
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
	sr.Results = new(Collection)
	sr.Status = response.Status
	
	parser.ParseResponseHeader(response, sr)
	
	if response.Status == 0 {
		err := parser.ParseResponse(response, sr)
		if err != nil {
			return nil, err
		}
		parser.ParseFacetCounts(response, sr)
		parser.ParseHighlighting(response, sr)
	} else {
		parser.ParseError(response, sr)
	}

	return sr, nil
}

func (parser *StandardResultParser) ParseResponseHeader(response *SelectResponse, sr *SolrResult) {
	if responseHeader, ok := response.Response["responseHeader"].(map[string]interface{}); ok {
		sr.ResponseHeader = responseHeader
	}
}

func (parser *StandardResultParser) ParseError(response *SelectResponse, sr *SolrResult) {
	if error, ok := response.Response["error"]; ok {
		sr.Error = error.(map[string]interface{})
	}
}

// ParseResponse will assign result and build sr.docs if there is a response.
// If there is no response or grouped property in response it will return error
func (parser *StandardResultParser) ParseResponse(response *SelectResponse, sr *SolrResult) (error) {
	var err error
	if resp, ok := response.Response["response"].(map[string]interface{}); ok {
		sr.Results.NumFound = int(resp["numFound"].(float64))
		sr.Results.Start = int(resp["start"].(float64))
		if docs, ok := resp["docs"].([]interface{}); ok {
			sr.Results.Docs = make([]Document, len(docs))
			for i, v := range docs {
				sr.Results.Docs[i] = Document(v.(map[string]interface{}))
			}
		}
	} else if grouped, ok := response.Response["grouped"].(map[string]interface{}); ok {
		sr.Grouped = grouped
	} else {
		err = fmt.Errorf(`Standard parser can only parse solr response with response object,
					ie response.response and response.response.docs. Or grouped response
					Please use other parser or implement your own parser`)
	}
	
	return err
}

// ParseFacetCounts will assign facet_counts to sr if there is one.
// No modification done here
func (parser *StandardResultParser) ParseFacetCounts(response *SelectResponse, sr *SolrResult) {
	if facetCounts, ok := response.Response["facet_counts"]; ok {
		sr.FacetCounts = facetCounts.(map[string]interface{})
	}
}

// ParseHighlighting will assign highlighting to sr if there is one.
// No modification done here
func (parser *StandardResultParser) ParseHighlighting(response *SelectResponse, sr *SolrResult) {
	if highlighting, ok := response.Response["highlighting"]; ok {
		sr.Highlighting = highlighting.(map[string]interface{})
	}
}
