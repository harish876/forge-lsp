package configstore

import (
	"fmt"
	"strings"

	"github.com/harish876/forge-lsp/utils"
	ini "github.com/harish876/go-tree-sitter/ini"
	"github.com/harish876/go-tree-sitter/python"
)

func GetSectionNameFromLsHint(sourceCode []byte) (string, error) {
	lang := python.GetLanguage()
	query := []byte(`
	(
		module (
        	(comment) @comment
            (#match? @comment "ls-hint-section_name")
        )
	)	
	`)
	q, err := GetQueryCursor(lang, sourceCode, query)
	if q.Node.HasError() {
		return "", fmt.Errorf("Syntax Tree has errors")
	}

	if err != nil {
		return "", err
	}
	q.Cursor.Exec(q.Query, q.Node)

	for {
		m, ok := q.Cursor.NextMatch()
		if !ok {
			break
		}
		m = q.Cursor.FilterPredicates(m, q.SourceCode)
		for _, c := range m.Captures {
			if c.Node.Type() == "comment" {
				sectionName := c.Node.Content(sourceCode)
				res := strings.Split(sectionName, ":")
				if len(res) == 2 {
					return strings.Trim(res[1], " "), nil
				}
			}
		}
	}
	return "", nil
}

func GetSettingNameByLine(sourceCode []byte, line int) ([]string, error) {
	var result []string
	lang := python.GetLanguage()
	query := []byte(`
	(
		(call
		  function: (attribute
			object: (identifier) @object
			attribute: (identifier) @method)
		  arguments: (argument_list
			(string
			  (string_content) @string_content
			)))
		(#eq? @object "config")
		(#eq? @method "get")
	  )
	`)
	q, err := GetQueryCursor(lang, sourceCode, query)
	if q.Node.HasError() {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	q.Cursor.Exec(q.Query, q.Node)

	for {
		m, ok := q.Cursor.NextMatch()
		if !ok {
			break
		}
		m = q.Cursor.FilterPredicates(m, q.SourceCode)
		for _, c := range m.Captures {
			if c.Node.Type() == "string_content" {
				settingVal := c.Node.Content(sourceCode)
				if uint32(line) == c.Node.StartPoint().Row {
					result = append(result, settingVal)
				}
			}
		}
	}
	return result, nil
}

func (store *ConfigStore) UpdateSections(sourceCode []byte) error {
	logger := utils.GetLogger()
	lang := ini.GetLanguage()
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

	q, err := GetQueryCursor(lang, sourceCode, query)
	if q.Node.HasError() {
		logger.Debug("Syntax Tree has errors")
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

		q, err := GetQueryCursor(lang, sourceCode, query)
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
			var row, lcol, rcol int
			for _, c := range m.Captures {
				switch q.Query.CaptureNameForId(c.Index) {
				case "name":
					name = string(sourceCode[c.Node.StartByte():c.Node.EndByte()])
					row = int(c.Node.StartPoint().Row)
					lcol = int(c.Node.StartPoint().Column)
					rcol = int(c.Node.EndPoint().Column)
				case "value":
					value = string(sourceCode[c.Node.StartByte():c.Node.EndByte()])
				}
			}
			if len(name) != 0 && len(value) != 0 {
				settingsMap[name] = Setting{
					Key:   name,
					Value: value,
					Metadata: SettingMetadata{
						Row:      row,
						LeftCol:  lcol,
						RightCol: rcol,
					},
				}
			}
		}
		if value, ok := store.Sections[key]; ok {
			value.Settings = settingsMap
			store.Sections[key] = value
		}
	}
	return nil
}
