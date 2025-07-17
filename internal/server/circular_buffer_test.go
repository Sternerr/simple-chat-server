package server

import (
	"testing"

	. "github.com/sternerr/termtalk/pkg/models"
)

func TestCircularBufferAdd(t *testing.T) {
	cb := NewCircularBuffer(3)

	cb.Add(Message{Type: MessageTypeChat, Message: "A"})
	if cb.Buffer[cb.R].Message != "A" {
		t.Errorf("Expected 'A', got '%s'", cb.Buffer[cb.R].Message)
	}

	cb.Add(Message{Type: MessageTypeChat, Message: "B"})
	cb.Add(Message{Type: MessageTypeChat, Message: "C"})

	if cb.Buffer[0].Message != "A" || cb.Buffer[1].Message != "B" || cb.Buffer[2].Message != "C" {
		t.Errorf("Buffer mismatch: %+v", cb.Buffer)
	}

	cb.Add(Message{Type: MessageTypeChat, Message: "D"})

	if cb.Buffer[0].Message != "D" {
		t.Errorf("Expected 'D' at buffer[0], got '%s'", cb.Buffer[0].Message)
	}

	if cb.F != 1 {
		t.Errorf("Expected F to be 1, got %d", cb.F)
	}
	if cb.R != 0 {
		t.Errorf("Expected R to be 0, got %d", cb.R)
	}
}

func testCircularBufferGetAll(t *testing.T) {
	cb := NewCircularBuffer(3)

	cb.Add(Message{Type: MessageTypeChat, Message: "A"})
	cb.Add(Message{Type: MessageTypeChat, Message: "B"})
	cb.Add(Message{Type: MessageTypeChat, Message: "C"})

	result := cb.GetAll()
	if len(result) != 3 || result[0].Message != "A" || result[1].Message != "B" || result[2].Message != "C" {
		t.Errorf("GetAll mismatch: %+v", result)
	}

	cb.Add(Message{Type: MessageTypeChat, Message: "D"})
	result = cb.GetAll()
	if len(result) != 3 || result[0].Message != "B" || result[1].Message != "C" || result[2].Message != "D" {
		t.Errorf("GetAll after overwrite: %+v", result)
	}
}
