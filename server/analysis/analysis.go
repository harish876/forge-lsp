package analysis

import (
	"context"
	"fmt"
	"strings"

	configstore "github.com/harish876/forge-lsp/config_store"
	"github.com/harish876/forge-lsp/pkg/lsp"
	"github.com/harish876/forge-lsp/utils"
	sitter "github.com/harish876/go-tree-sitter"
	"github.com/harish876/go-tree-sitter/python"
)

var (
	logger = utils.GetLogger("/var/logs/forge-lsp.vscode.log")
)

type Document struct {
	Content string
	Node    *sitter.Node
}

type State struct {
	Documents map[string]*Document
}

func NewState() State {
	return State{
		Documents: map[string]*Document{},
	}
}

func (s *State) OpenDocument(uri, text string) {
	node, _ := sitter.ParseCtx(context.Background(), []byte(text), python.GetLanguage())
	s.Documents[uri] = &Document{
		Content: text,
		Node:    node,
	}
}

func (s *State) UpdateDocument(uri, contentChange string, store *configstore.ConfigStore) {
	if strings.Contains(uri, ".ini") {
		store.UpdateSections([]byte(contentChange))
	}
	s.Documents[uri].Content = contentChange
	newTree, err := sitter.ParseCtx(context.Background(), []byte(contentChange), python.GetLanguage())
	if err != nil {
		s.Documents[uri].Node = newTree
	}
}

func (s *State) Hover(id int, uri string, position int) lsp.HoverResponse {
	document := s.Documents[uri]
	return lsp.HoverResponse{
		Response: lsp.Response{
			ID:  &id,
			RPC: "2.0",
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("Document: %s  Characters: %d", uri, len(document.Content)),
		},
	}

}

func (s *State) Definition(id int, uri string, line int, store *configstore.ConfigStore) lsp.DefinitionResponse {
	document := s.Documents[uri]
	section := utils.GetSectionNameFromUri(uri)
	settingNameFromCode := configstore.GetSettingNameByLine([]byte(document.Content), line)
	var setting configstore.Setting

	if len(settingNameFromCode) > 0 {
		if value, ok := store.Sections[section]; ok {
			if setting, ok = value.Settings[settingNameFromCode[0]]; !ok {
				logger.Printf("unable to find setting %s", settingNameFromCode[0])
			}
		}
	}
	return lsp.DefinitionResponse{
		Response: lsp.Response{
			ID:  &id,
			RPC: "2.0",
		},
		Result: lsp.Location{
			Uri: store.BasePath,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      setting.Metadata.Row,
					Character: setting.Metadata.LeftCol,
				},
				End: lsp.Position{
					Line:      setting.Metadata.Row,
					Character: setting.Metadata.RightCol,
				},
			},
		},
	}
}
