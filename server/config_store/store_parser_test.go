package configstore

import (
	"fmt"
	"testing"
)

func TestStoreParser(t *testing.T) {
	store := NewConfigStore()
	sourceCode, err := store.OpenConfigFile("/Users/harishgokul/forge-lsp/server/config_store/settings.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = store.GetSections(sourceCode)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(store.Sections)
}
