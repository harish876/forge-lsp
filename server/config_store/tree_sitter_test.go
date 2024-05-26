package configstore

import (
	"context"
	"fmt"
	"testing"

	sitter "github.com/harish876/go-tree-sitter"
	ini "github.com/harish876/go-tree-sitter/ini"
)

func TestTreeSitterEdit(t *testing.T) {
	parser := sitter.NewParser()
	parser.SetLanguage(ini.GetLanguage())
	sourceCode := []byte("let a = 1")
	oldTree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
	// change 1 -> true
	newText := []byte("let a = true")
	oldTree.Edit(sitter.EditInput{
		StartIndex:  8,
		OldEndIndex: 9,
		NewEndIndex: 12,
		StartPoint: sitter.Point{
			Row:    0,
			Column: 8,
		},
		OldEndPoint: sitter.Point{
			Row:    0,
			Column: 9,
		},
		NewEndPoint: sitter.Point{
			Row:    0,
			Column: 12,
		},
	})

	// generate new tree
	newTree, _ := parser.ParseCtx(context.Background(), oldTree, newText)
	fmt.Println(newTree)
}
