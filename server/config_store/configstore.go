package configstore

import (
	"context"
	"fmt"
	"os"

	"github.com/harish876/forge-lsp/utils"
	sitter "github.com/harish876/go-tree-sitter"
	ini "github.com/harish876/go-tree-sitter/ini"
)

var (
	SECTION_NODE_PARENT = "section_name"
	SECTION_NODE_TYPE   = "text"
	logger              = utils.GetLogger("/Users/harishgokul/forge-lsp/server/log.txt")
)

type Setting struct {
	Key   string
	Value string
}

type Section struct {
	Settings map[string]string
}

func NewSection() Section {
	return Section{
		Settings: make(map[string]string, 0),
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
	Tree       *sitter.Tree
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

func GetQueryCursor(sourceCode []byte, query []byte) (*QueryExecutionParams, error) {
	lang := ini.GetLanguage()
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

func (store *ConfigStore) UpdateSections(sourceCode []byte) error {
	query := []byte(`
	(
		document(
		 section(
		  section_name (
			  text
		  ) @sectionName
		 )
		)
	  )
	`)

	q, err := GetQueryCursor(sourceCode, query)
	if q.Node.HasError() {
		logger.Println("Syntax Tree has errors")
	}

	if err != nil {
		fmt.Println(err)
	}
	q.Cursor.Exec(q.Query, q.Node)

	for {
		m, ok := q.Cursor.NextMatch()
		if !ok {
			break
		}
		m = q.Cursor.FilterPredicates(m, q.SourceCode)
		for _, c := range m.Captures {
			store.Sections[c.Node.Content(sourceCode)] = NewSection()
		}
	}

	for key := range store.Sections {
		settingsMap := store.Sections[key].Settings
		_ = settingsMap
		query := []byte(fmt.Sprintf(`
		(
			document(
			 section(
			  section_name (
				  text
			  ) @sectionName
			  (#match? @sectionName %s)
			 )
			 (setting
				setting_name: (setting_name) @name
				setting_value: (setting_value) @value
			 ) @setting
		   )
		)
		`, key))

		q, err := GetQueryCursor(sourceCode, query)
		// if q.Node.HasError() {
		// 	logger.Println("Syntax Tree has errors")
		// 	continue
		// }
		if err != nil {
			return err
		}
		q.Cursor.Exec(q.Query, q.Node)

		for {
			m, ok := q.Cursor.NextMatch()
			if !ok {
				break
			}
			m = q.Cursor.FilterPredicates(m, q.SourceCode)
			var name, value string
			for _, c := range m.Captures {
				switch q.Query.CaptureNameForId(c.Index) {
				case "name":
					name = string(sourceCode[c.Node.StartByte():c.Node.EndByte()])
				case "value":
					value = string(sourceCode[c.Node.StartByte():c.Node.EndByte()])
				}
			}
			if len(name) != 0 && len(value) != 0 {
				settingsMap[name] = value
			}
		}
		if value, ok := store.Sections[key]; ok {
			value.Settings = settingsMap
			store.Sections[key] = value
		}
	}
	return nil
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
		for key, _ := range value.Settings {
			result = append(result, key)
		}
	}
	return result
}

func (store *ConfigStore) ListSettings(section string) []Setting {
	var result []Setting
	if value, ok := store.Sections[section]; ok {
		for key, value := range value.Settings {
			//TODO: get the key as well
			result = append(result, Setting{Key: key, Value: value})
		}
	}
	return result
}