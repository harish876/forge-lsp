package configstore

import (
	"fmt"
	"testing"
)

// func TestStoreParser(t *testing.T) {
// 	store := NewConfigStore()
// 	sourceCode, err := store.OpenConfigFile("/Users/harishgokul/forge-lsp/server/config_store/settings.ini")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	err = store.UpdateSections(sourceCode)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println(store.Sections)
// }

func TestGetSettingFromCapture(t *testing.T) {
	sourceCode := []byte(`
	import logging
	import pandas as pd
	from db.mssql_service import MssqlClient
	from jobs.job_interface import ETLJob
	from sqlalchemy import URL
	import os
	from datetime import datetime,timedelta
	
	#ls-hint-section_name: extract_odin
	
	class OdinDataExtractJob(ETLJob):
	
		def __init__(self,config):        
			try:
				self.create_client()
			except Exception as e:
				logging.error(e)
			finally:
				self.chunk  = config.get('chunk')
				self.query = config.get('query')
				config.get('type')
				
				if config.get('delta_days') is not None:
					self.delta = int(config.get('delta_days'))
				time_delta = datetime.now() - timedelta(days=self.delta)
				self.unix_timestamp_end = int(datetime.now().timestamp())
				self.unix_timestamp_start = int(time_delta.timestamp())
				self.params = (self.unix_timestamp_start,self.unix_timestamp_end)
	
		
		def create_client(self):
			connection_url = URL.create(
				"mssql+pyodbc",
				host = os.environ.get("ODIN_SERVER", None),
				database = os.environ.get("ODIN_DATABASE", None),
				username = os.environ.get("ODIN_USERNAME", None),
				password = os.environ.get("ODIN_PASSWORD", None),
				port = os.environ.get("ODIN_PORT", 1433),
				query= {
					"driver":"ODBC Driver 13 for SQL Server",
					"TrustServerCertificate":"Yes"
				})
			
			#init the mssql client here for odin
			mssql_client = MssqlClient(config={
				"db_url":connection_url,
				"stream_results":True,
				"trust_server_certificate":"yes",
				"driver":"ODBC Driver 13 for SQL Server"
			})
			try:
				self.client = mssql_client.get_client("mssql+pyodbc")
			except Exception as e:
				raise e
		
		def execute(self):
			try:
				for chunk_dataframe in pd.read_sql(self.query,self.client,chunksize = int(self.chunk),params=self.params):
					if len(chunk_dataframe) == 0:
						logging.info("No Data Available")
										
					logging.info("Dataframe: ",chunk_dataframe)
					self.set_data_context(chunk_dataframe)
					self.next()
	
			except Exception as e:
				logging.error(e)
				raise e	
	`)
	res, err := GetSettingNameByLine(sourceCode, 4)
	fmt.Println("Result", res)
	if err != nil {
		t.Fatal(err)
	}
}
