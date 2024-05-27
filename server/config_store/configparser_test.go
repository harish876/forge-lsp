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
	class ExtractCsvJob(ETLJob):

    def __init__(self, config):
        super()
        self.__filename = config.get("directory")

    def execute(self, data=None):
        try:
            if self.__filename is None:
                return

            data = pd.read_csv(self.__filename)
            self.set_data_context(data.head())
            self.next()

        except Exception as e:
            logging.error(e)
            raise e

        if data.empty:
            return
	`)
	GetSettingNameByLine(sourceCode, 4)
}
