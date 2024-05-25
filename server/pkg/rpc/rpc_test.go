package rpc_test

import (
	"testing"

	"github.com/harish876/forge-lsp/pkg/rpc"
)

type EncodingExample struct {
	Method string
}

func TestEncoding(t *testing.T) {
	expected := "Content-Length: 15\r\n\r\n{\"method\":\"hi\"}"
	actual := rpc.EncodeMessage(EncodingExample{Method: "hi"})

	if expected == actual {
		t.Fatalf("Expected: %s,  Actual: %s", expected, actual)
	}
}

func TestDecoding(t *testing.T) {
	incomingMessage := "Content-Length: 15\r\n\r\n{\"method\":\"hi\"}"
	method, content, err := rpc.DecodeMessage([]byte(incomingMessage))

	if err != nil {
		t.Fatalf("unable to decode message: %v", err)
	}

	contentLen := len(content)
	if contentLen != 15 {
		t.Fatalf("Expected 15 but got content length: %d", contentLen)
	}

	if method != "hi" {
		t.Fatalf("Expected hi but got %s", method)
	}
}
