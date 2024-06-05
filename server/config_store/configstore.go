package configstore

import (
	"context"
	"fmt"
	"os"

	sitter "github.com/harish876/go-tree-sitter"
)

var (
	SECTION_NODE_PARENT = "section_name"
	SECTION_NODE_TYPE   = "text"
)

type SettingMetadata struct {
	Row      int
	LeftCol  int
	RightCol int
}

type Setting struct {
	Key      string
	Value    string
	Metadata SettingMetadata
}

type Section struct {
	Settings map[string]Setting
}

func NewSection() Section {
	return Section{
		Settings: make(map[string]Setting, 0),
	}
}

type ConfigStore struct {
	BasePath string
	Sections map[string]Section
}

func NewConfigStore() *ConfigStore {
	return &ConfigStore{
		Sections: make(map[string]Section, 0),
	}
}

type QueryExecutionParams struct {
	Cursor     *sitter.QueryCursor
	Query      *sitter.Query
	Node       *sitter.Node
	SourceCode []byte
}

func NewQueryExecutionParams(cursor *sitter.QueryCursor, query *sitter.Query, node *sitter.Node, sourceCode []byte) *QueryExecutionParams {
	return &QueryExecutionParams{
		Cursor:     cursor,
		Query:      query,
		Node:       node,
		SourceCode: sourceCode,
	}
}

func GetQueryCursor(lang *sitter.Language, sourceCode []byte, query []byte) (*QueryExecutionParams, error) {
	node, _ := sitter.ParseCtx(context.Background(), sourceCode, lang)

	sitterQuery, _ := sitter.NewQuery(query, lang)
	queryCursor := sitter.NewQueryCursor()

	return NewQueryExecutionParams(queryCursor, sitterQuery, node, sourceCode), nil
}

func (store *ConfigStore) OpenConfigFile(filePath string) ([]byte, error) {
	store.BasePath = filePath
	if len(store.BasePath) == 0 {
		return nil, fmt.Errorf("no config file found at %s", store.BasePath)
	}
	codeBuf := make([]byte, 10*1024*1024)
	file, err := os.Open(store.BasePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	n, err := file.Read(codeBuf)
	if err != nil {
		return nil, err
	}
	sourceCode := codeBuf[:n]

	return sourceCode, nil
}

func (store *ConfigStore) ListSections() []string {
	var result []string
	for key := range store.Sections {
		result = append(result, key)
	}
	return result
}

func (store *ConfigStore) ListAllSettings() []string {
	var result []string
	for _, value := range store.Sections {
		for key := range value.Settings {
			result = append(result, key)
		}
	}
	return result
}

func (store *ConfigStore) ListSettings(section string) []Setting {
	var result []Setting
	if value, ok := store.Sections[section]; ok {
		for _, value := range value.Settings {
			result = append(result, value)
		}
	}
	return result
}
