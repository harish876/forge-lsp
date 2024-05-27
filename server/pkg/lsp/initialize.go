package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
	RootUri    string      `json:"rootUri"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync   int                `json:"textDocumentSync"`
	PositionEncoding   string             `json:"positionEncoding,omitempty"`
	HoverProvider      bool               `json:"hoverProvider"`
	DefinitionProvider bool               `json:"definitionProvider"`
	CompletionProvider CompletionProvider `json:"completionProvider,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type CompletionProvider struct {
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			ID:  &id,
			RPC: "2.0",
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   1,
				PositionEncoding:   "utf-16",
				HoverProvider:      true,
				DefinitionProvider: true,
				CompletionProvider: CompletionProvider{},
			},
			ServerInfo: ServerInfo{
				Name:    "forge-lsp",
				Version: "0.0.1",
			},
		},
	}
}
