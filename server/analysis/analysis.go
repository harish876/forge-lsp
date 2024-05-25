package analysis

import (
	"fmt"

	"github.com/harish876/forge-lsp/pkg/lsp"
)

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: map[string]string{},
	}
}

func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) UpdateDocument(uri, contentChange string) {
	s.Documents[uri] = contentChange
}

func (s *State) Hover(id int, uri string, position int) lsp.HoverResponse {
	document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			ID:  &id,
			RPC: "2.0",
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("Document: %s  Characters: %d", uri, len(document)),
		},
	}

}
