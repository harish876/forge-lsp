package lsp

import (
	"github.com/harish876/forge-lsp/config_store"
)

type TextDocumentItem struct {
	/**
	 * The text document's URI.
	 */
	URI string `json:"uri"`

	/**
	 * The text document's language identifier.
	 */
	LanguageID string `json:"languageId"`

	/**
	 * The version number of this document (it will increase after each
	 * change, including undo/redo).
	 */
	Version int `json:"version"`

	/**
	 * The content of the opened text document.
	 */
	Text string `json:"string"`
}

type TextDocumentIdentifier struct {
	/**
	 * The text document's URI.
	 */
	URI string `json:"uri"`
}

type DidOpenTextDocumentNotification struct {
	Notification
	Params DidOpenTextDocumentParams `json:"params"`
}

type DidOpenTextDocumentParams struct {
	/**
	 * The document that was opened.
	 */
	TextDocument TextDocumentItem `json:"textDocument"`
}

type DidChangeTextDocumentNotification struct {
	Notification
	Params DidChangeTextDocumentParams `json:"params"`
}

type DidChangeTextDocumentParams struct {
	/**
	 * The document that did change. The version number points
	 * to the version after all provided content changes have
	 * been applied.
	 */
	TextDocument VersionedTextDocumentIdentifier `json:"textDocument"`

	/**
	 * The actual content changes. The content changes describe single state
	 * changes to the document. So if there are two content changes c1 (at
	 * array index 0) and c2 (at array index 1) for a document in state S then
	 * c1 moves the document from S to S' and c2 from S' to S''. So c1 is
	 * computed on the state S and c2 is computed on the state S'.
	 *
	 * To mirror the content of a document using change events use the following
	 * approach:
	 * - start with the same initial content
	 * - apply the 'textDocument/didChange' notifications in the order you
	 *   receive them.
	 * - apply the `TextDocumentContentChangeEvent`s in a single notification
	 *   in the order you receive them.
	 */
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

type TextDocumentContentChangeEvent struct {
	//The new text of the whole document.
	Text string `json:"text"`
}

type HoverRequest struct {
	Request
	Params HoverRequestParams `json:"params"`
}

type HoverRequestParams struct {
	TextDocumentPositionParams
}

type TextDocumentPositionParams struct {
	/**
	 * The text document.
	 */
	TextDocument TextDocumentIdentifier `json:"textDocument"`

	/**
	 * The position inside the text document.
	 */
	Position Position `json:"position"`
}

type Position struct {
	/**
	 * Line position in a document (zero-based).
	 */
	Line int `json:"line"`

	/**
	 * Character offset on a line in a document (zero-based). The meaning of this
	 * offset is determined by the negotiated `PositionEncodingKind`.
	 *
	 * If the character value is greater than the line length it defaults back
	 * to the line length.
	 */
	Character int `json:"character"`
}

type HoverResponse struct {
	Response
	Result HoverResult `json:"result"`
}

type HoverResult struct {
	Contents string `json:"contents"`
}

type TextDocumentCompletionRequest struct {
	Request
	Params TextDocumentCompletionParams `json:"params"`
}

type TextDocumentCompletionParams struct {
	TextDocumentPositionParams
	Context CompletionContext `json:"context,omitempty"`
}

type CompletionContext struct {
	/**
	 * How the completion was triggered.
	 */
	TriggerKind int `json:"triggerKind"`

	/**
	 * The trigger character (a single character) that has trigger code
	 * complete. Is undefined if
	 * `triggerKind !== CompletionTriggerKind.TriggerCharacter`
	 */
	TriggerCharacter string `json:"triggerCharacter,omitempty"`
}

type TextDocumentCompletionResponse struct {
	Response
	Result CompletionList `json:"result"`
}

type CompletionList struct {
	/**
	 * This list is not complete. Further typing should result in recomputing
	 * this list.
	 *
	 * Recomputed lists have all their items replaced (not appended) in the
	 * incomplete completion sessions.
	 */
	IsIncomplete bool `json:"isIncomplete"`

	/**
	 * In many cases the items of an actual completion result share the same
	 * value for properties like `commitCharacters` or the range of a text
	 * edit. A completion list can therefore define item defaults which will
	 * be used if a completion item itself doesn't specify the value.
	 *
	 * If a completion list specifies a default value and a completion item
	 * also specifies a corresponding value the one from the item is used.
	 *
	 * Servers are only allowed to return default values if the client
	 * signals support for this via the `completionList.itemDefaults`
	 * capability.
	 *
	 * @since 3.17.0
	 */
	/**
	 * The completion items.
	 */
	Items []CompletionItem `json:"items"`
}

type CompletionItem struct {
	/**
	 * The label of this completion item.
	 *
	 * The label property is also by default the text that
	 * is inserted when selecting this completion.
	 *
	 * If label details are provided the label itself should
	 * be an unqualified name of the completion item.
	 */
	Label string `json:"label"`
}

func NewTextDocumentCompletionResponse(id int, store *configstore.ConfigStore) TextDocumentCompletionResponse {
	sectionList := store.ListSettings()
	var items []CompletionItem
	for _, section := range sectionList {
		items = append(items, CompletionItem{
			Label: section,
		})
	}
	return TextDocumentCompletionResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: CompletionList{
			IsIncomplete: true,
			Items:        items,
		},
	}
}
