package lsp

type Request struct {
	RPC    string `json:"rpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`

	//Todo: Later
	//Params ...
}

type Response struct {
	RPC string `json:"rpc"`
	ID  *int   `json:"id,omitempty"`

	//Result ...
	//Error ...
}

type Notification struct {
	RPC    string `json:"rpc"`
	Method string `json:"method"`
}
