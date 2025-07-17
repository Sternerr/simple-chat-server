package protocol_test

import (
	"testing"

	"github.com/sternerr/termtalk/internal/protocol"
	. "github.com/sternerr/termtalk/pkg/models"
)

func TestEncodeHandshake(t *testing.T) {
	expected := `{"type":"handshake","from":"username"}`

	actual, err := protocol.EncodeMessage(Message{
		Type: "handshake",
		From: "username",
	})

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if string(actual) != expected {
		t.Fatalf("Expected: %q, Actual: %q", expected, string(actual))
	}
}

func TestDecodeHandshake(t *testing.T) {
	input := []byte("{\"type\":\"handshake\",\"from\":\"username\"}")

	expected := &Message{
		Type: "handshake",
		From: "username",
	}

	actual, err := protocol.DecodeMessage(input)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if actual.Type != expected.Type || actual.From != expected.From {
		t.Fatalf("Decoded message doesn't match.\nExpected: %+v\nActual: %+v", expected, actual)
	}
}

func TestEncodeChat(t *testing.T) {
	expected := `{"type":"chat","from":"username","message":"Hello World"}`

	actual, err := protocol.EncodeMessage(Message{
		Type: "chat",
		From: "username",
		Message: "Hello World",
	})

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if string(actual) != expected {
		t.Fatalf("Expected: %q, Actual: %q", expected, string(actual))
	}
}

func TestDecodeChat(t *testing.T) {
	input := []byte("{\"type\":\"chat\",\"from\":\"username\",\"message\":\"Hello World\"}")

	expected := Message{
		Type: "chat",
		From: "username",
		Message: "Hello World",
	}

	actual, err := protocol.DecodeMessage(input)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if actual.Type != expected.Type || actual.From != expected.From || actual.Message != expected.Message {
		t.Fatalf("Decoded message doesn't match.\nExpected: %+v\nActual: %+v", expected, actual)
	}
}
