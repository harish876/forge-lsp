package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/harish876/forge-lsp/analysis"
	configstore "github.com/harish876/forge-lsp/config_store"
	"github.com/harish876/forge-lsp/pkg/lsp"
	rpc "github.com/harish876/forge-lsp/pkg/rpc"
	"github.com/harish876/forge-lsp/utils"
)

type Args struct {
	LogFile  string
	LogLevel string
}

func ParseArgs() (*Args, error) {
	logFile := flag.String("file", "forge-lsp.vscode.log", "Path to Log File")
	logLevel := flag.String("level", "INFO", "Log Level [ INFO | ERROR | DEBUG ]")
	logStdio := flag.Bool("stdio", true, "Default flag. IDK")
	_ = logStdio

	flag.Parse()

	return &Args{
		LogFile:  *logFile,
		LogLevel: *logLevel,
	}, nil
}

func main() {
	args, err := ParseArgs()
	if err != nil {
		log.Fatal(err)
	}
	logger := utils.NewLogger(args.LogFile, args.LogLevel)
	logger.Debug(fmt.Sprintf("LSP Server Started %s", args.LogLevel))
	configStore := configstore.NewConfigStore()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	state := analysis.NewState()
	writer := os.Stdout

	for {
		for scanner.Scan() {
			msg := scanner.Bytes()
			method, content, err := rpc.DecodeMessage(msg)
			if err != nil {
				logger.Error("We gots some error: %v", err)
			}
			handlerMessage(logger, method, content, state, writer, configStore)
		}
	}

}

func handlerMessage(logger *slog.Logger, method string, content []byte, state analysis.State, writer io.Writer, store *configstore.ConfigStore) {
	logger.Debug(fmt.Sprintf("Received method: %s", method))
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Error("Could Not Unmarshal initialize request request %v", err)
		}
		logger.Info(
			"Initialize Request",
			"Name",
			request.Params.ClientInfo.Name,
			"Version",
			request.Params.ClientInfo.Name,
			"RootURI",
			request.Params.RootUri,
		)

		go initStore(request, store)

		msg := lsp.NewInitializeResponse(request.ID)
		reply := writeResponse(writer, msg)

		logger.Info("Sent the message", "reply", reply)

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Error("Could Not Unmarshal initialize textDocument/didOpen request %v", err)
		}
		logger.Info("textDocument/didOpen", "URI",
			request.Params.TextDocument.URI,
		)

		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)

	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Error("Could Not Unmarshal initialize textDocument/didChange request %v", err)
		}
		logger.Info("textDocument/didChange", "URI",
			request.Params.TextDocument.URI)

		for _, contentChange := range request.Params.ContentChanges {
			state.UpdateDocument(request.Params.TextDocument.URI, contentChange.Text, store)
		}

	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Error("Could Not Unmarshal initialize textDocument/hover request %v", err)
		}
		logger.Info("textDocument/hover",
			"URI",
			request.Params.TextDocument.URI,
			"Line",
			request.Params.Position.Character,
			"Character",
			request.Params.Position.Line,
		)

		msg := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position.Line)
		reply := writeResponse(writer, msg)
		logger.Info("Sent the reply for textDocumen/hover", "textDocumen/hover", reply)

	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Error("Could Not Unmarshal initialize textDocument/definition request %v", err)
		}
		logger.Info("textDocument/definition",
			"URI",
			request.Params.TextDocument.URI,
			"Line",
			request.Params.Position.Line,
			"Character",
			request.Params.TextDocumentPositionParams.Position.Character,
		)

		msg := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position.Line, store)
		reply := writeResponse(writer, msg)
		logger.Info("Sent the reply for textDocumen/definition", "textDocumen/definition", reply)

	case "textDocument/completion":
		var request lsp.TextDocumentCompletionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Error("Could Not Unmarshal initialize textDocument/completion request %v", err)
		}
		logger.Info("textDocument/completion",
			"TriggerCharacter",
			request.Params.Context.TriggerCharacter,
		)

		msg := lsp.NewTextDocumentCompletionResponse(request.ID, request.Params.TextDocument.URI, store)
		reply := writeResponse(writer, msg)
		logger.Info("Sent the reply for textDocumen/completion", "textDocumen/completion", reply)
	}
}

func writeResponse(writer io.Writer, msg any) string {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
	return reply
}

func initStore(request lsp.InitializeRequest, store *configstore.ConfigStore) {
	logger := utils.GetLogger()
	uri := strings.TrimPrefix(request.Params.RootUri, "file:")
	path := filepath.Join(uri, "config", "settings.local.ini")
	logger.Debug(path)
	sourceCode, err := store.OpenConfigFile(path)
	if err != nil {
		logger.Error("Error at OpenConfigFile", err)
	}
	store.UpdateSections(sourceCode)
	logger.Info("Config Section", "Section", store.Sections)
}
