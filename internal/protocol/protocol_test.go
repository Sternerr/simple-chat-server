package protocol_test

import (
	"testing"
	"github.com/sternerr/termtalk/internal/protocol"
)

func TestDecodeHandshake(t *testing.T) {
	expected := "username"
	actual, err := protocol.DecodeHandshake([]byte("Handshake\r\n\r\nusername"))
	if err != nil {
		t.Fatalf("Expected: %s, Actual: %s", expected, err.Error())
	}

	if actual != expected {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}
