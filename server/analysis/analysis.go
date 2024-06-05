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

func (s *State) Hover(id int, uri string, line int, store *configstore.ConfigStore) lsp.HoverResponse {
	document := s.Documents[uri]
	logger := utils.GetLogger()

	section := utils.GetSectionNameFromUri(uri)

	settingNameFromCode, err := configstore.GetSettingNameByLine([]byte(document.Content), line)
	if err != nil {
		logger.Error("Definition", "Setting Name From Code Error", err)
	}
	var setting configstore.Setting

	if _, ok := store.Sections[section]; !ok {
		section, _ = configstore.GetSectionNameFromLsHint([]byte(document.Content)) //
	}

	logger.Debug("Definition", "Section Name From Ls Hint", section)
	logger.Debug("Definition", "Setting", settingNameFromCode)
	if len(settingNameFromCode) > 0 {
		if value, ok := store.Sections[section]; ok {
			if setting, ok = value.Settings[settingNameFromCode[0]]; !ok {
				logger.Debug(fmt.Sprintf("unable to find setting %s", settingNameFromCode[0]))
			}
		}
	}
	var content string
	if len(section) == 0 {
		content = "No such section present"
	} else if len(setting.Key) == 0 {
		content = fmt.Sprintf("No such value key present under %s", section)
	} else {
		content = fmt.Sprintf("Section - %s Key - %s  Value- %s", section, setting.Key, setting.Value)
	}

	return lsp.HoverResponse{
		Response: lsp.Response{
			ID:  &id,
			RPC: "2.0",
		},
		Result: lsp.HoverResult{
			Contents: content,
		},
	}

}

func (s *State) Definition(id int, uri string, line int, store *configstore.ConfigStore) lsp.DefinitionResponse {
	logger := utils.GetLogger()
	document := s.Documents[uri]

	section := utils.GetSectionNameFromUri(uri)

	settingNameFromCode, err := configstore.GetSettingNameByLine([]byte(document.Content), line)
	if err != nil {
		logger.Error("Definition", "Setting Name From Code Error", err)
	}
	var setting configstore.Setting

	if _, ok := store.Sections[section]; !ok {
		section, _ = configstore.GetSectionNameFromLsHint([]byte(document.Content)) //
	}

	logger.Debug("Definition", "Section Name From Ls Hint", section)
	logger.Debug("Definition", "Setting", settingNameFromCode)
	if len(settingNameFromCode) > 0 {
		if value, ok := store.Sections[section]; ok {
			if setting, ok = value.Settings[settingNameFromCode[0]]; !ok {
				logger.Debug(fmt.Sprintf("unable to find setting %s", settingNameFromCode[0]))
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

func (s *State) Completion(id int, uri string, store *configstore.ConfigStore) lsp.TextDocumentCompletionResponse {
	document := s.Documents[uri]
	section := utils.GetSectionNameFromUri(uri)
	sectionList := store.ListSettings(section)

	if len(sectionList) == 0 {
		section, _ = configstore.GetSectionNameFromLsHint([]byte(document.Content))
		sectionList = store.ListSettings(section)
	}

	var items []lsp.CompletionItem
	for _, section := range sectionList {
		items = append(items, lsp.CompletionItem{
			Label:  section.Key,
			Detail: fmt.Sprintf("%s = %s", section.Key, section.Value),
			Kind:   5,
		})
	}
	return lsp.TextDocumentCompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.CompletionList{
			IsIncomplete: false,
			Items:        items,
		},
	}
}
