package protocol_test

import (
	"testing"
	"github.com/sternerr/termtalk/internal/protocol"
	. "github.com/sternerr/termtalk/pkg/models"
)

func TestEncodeHandshake(t *testing.T) {
	expected := "Handshake\r\n\r\nusername"
	actual := protocol.EncodeHandshake(User{
		Username: "username",
	})

	if actual != expected {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}

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
