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
	err = store.UpdateSections(sourceCode)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(store.Sections)
}

func TestGetSettingFromCapture(t *testing.T) {
	sourceCode := []byte(`
from jobs.job_interface import ETLJob

class ExtractJsonJob(ETLJob):
	def __init__(self, config):
     	config.get("foo")
		config.get("file_name")
       	config.get("filename")

	def execute(self, data=None):
		self.set_data_context('foobar')
		self.next()
	`)
	GetSettingNameByLine(sourceCode, 4)
}
