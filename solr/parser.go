package solr

import (
	"fmt"
)

// ResultParser is interface for parsing result from response.
// The idea here is that application have possibility to parse.
// Or defined own parser with internal data structure to suite
// application's need
type ResultParser interface {
	Parse(response *SolrResponse) (*SolrResult, error)
}

type StandardResultParser struct {
}

func (parser *StandardResultParser) Parse(response *SolrResponse) (*SolrResult, error) {
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
		parser.ParseStats(response, sr)
		parser.ParseMoreLikeThis(response, sr)
	} else {
		parser.ParseError(response, sr)
	}

	return sr, nil
}

func (parser *StandardResultParser) ParseResponseHeader(response *SolrResponse, sr *SolrResult) {
	if responseHeader, ok := response.Response["responseHeader"].(map[string]interface{}); ok {
		sr.ResponseHeader = responseHeader
	}
}

func (parser *StandardResultParser) ParseError(response *SolrResponse, sr *SolrResult) {
	if err, ok := response.Response["error"].(map[string]interface{}); ok {
		sr.Error = err
	}
}

func ParseDocResponse(docResponse map[string]interface{}, collection *Collection) {
	collection.NumFound = int(docResponse["numFound"].(float64))
	collection.Start = int(docResponse["start"].(float64))
	if docs, ok := docResponse["docs"].([]interface{}); ok {
		collection.Docs = make([]Document, len(docs))
		for i, v := range docs {
			collection.Docs[i] = Document(v.(map[string]interface{}))
		}
	}
}

// ParseSolrResponse will assign result and build sr.docs if there is a response.
// If there is no response or grouped property in response it will return error
func (parser *StandardResultParser) ParseResponse(response *SolrResponse, sr *SolrResult) (err error) {
	if resp, ok := response.Response["response"].(map[string]interface{}); ok {
		ParseDocResponse(resp, sr.Results)
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
func (parser *StandardResultParser) ParseFacetCounts(response *SolrResponse, sr *SolrResult) {
	if facetCounts, ok := response.Response["facet_counts"].(map[string]interface{}); ok {
		sr.FacetCounts = facetCounts
	}
}

// ParseHighlighting will assign highlighting to sr if there is one.
// No modification done here
func (parser *StandardResultParser) ParseHighlighting(response *SolrResponse, sr *SolrResult) {
	if highlighting, ok := response.Response["highlighting"].(map[string]interface{}); ok {
		sr.Highlighting = highlighting
	}
}

// Parse stats if there is  in response
func (parser *StandardResultParser) ParseStats(response *SolrResponse, sr *SolrResult) {
	if stats, ok := response.Response["stats"].(map[string]interface{}); ok {
		sr.Stats = stats
	}
}

// Parse moreLikeThis if there is in response
func (parser *StandardResultParser) ParseMoreLikeThis(response *SolrResponse, sr *SolrResult) {
	if moreLikeThis, ok := response.Response["moreLikeThis"].(map[string]interface{}); ok {
		sr.MoreLikeThis = moreLikeThis
	}
}

type MltResultParser interface {
	Parse(response *SolrResponse) (*SolrMltResult, error)
}

type MoreLikeThisParser struct {
}

func (parser *MoreLikeThisParser) Parse(response *SolrResponse) (*SolrMltResult, error) {
	sr := &SolrMltResult{}
	sr.Results = new(Collection)
	sr.Match = new(Collection)
	sr.Status = response.Status
	
	if responseHeader, ok := response.Response["responseHeader"].(map[string]interface{}); ok {
		sr.ResponseHeader = responseHeader
	}
	
	if response.Status == 0 {
		if resp, ok := response.Response["response"].(map[string]interface{}); ok {
			ParseDocResponse(resp, sr.Results)
		}
		if match, ok := response.Response["match"].(map[string]interface{}); ok {
			ParseDocResponse(match, sr.Match)
		}
	} else {
		if err, ok := response.Response["error"].(map[string]interface{}); ok {
			sr.Error = err
		}
	}
	return sr, nil
}
