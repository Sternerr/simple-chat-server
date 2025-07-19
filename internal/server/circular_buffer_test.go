package server_test

import "github.com/sternerr/termtalk/internal/server"
import . "github.com/sternerr/termtalk/pkg/models"
import "testing"

func TestCircularBufferAddAndGetAll(t *testing.T) {
	cb := server.NewCircularBuffer(3)

	cb.Add(Message{Type: MessageTypeText, From: "client1", Message: "A"})
	cb.Add(Message{Type: MessageTypeText, From: "client2", Message: "B"})

	messages := cb.GetAll()
	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}

	if messages[0].Message != "A" || messages[1].Message != "B" {
		t.Errorf("Unexpected message order or content")
	}
}


func TestCircularBufferOverwrite(t *testing.T) {
	cb := server.NewCircularBuffer(2)
	
	cb.Add(Message{Type: MessageTypeText, From: "client1", Message: "A"})
	cb.Add(Message{Type: MessageTypeText, From: "client2", Message: "B"})
	cb.Add(Message{Type: MessageTypeText, From: "client3", Message: "C"})

	messages := cb.GetAll()
	if len(messages) != 2 {
		t.Errorf("Expected 2 messages after overwrite, got %d", len(messages))
	}

	if messages[0].Message != "B" || messages[1].Message != "C" {
		t.Errorf("Unexpected message order or content after overwrite")
	}
}

func TestCircularBufferExactCapacity(t *testing.T) {
	cb := server.NewCircularBuffer(3)

	cb.Add(Message{Type: MessageTypeText, From: "client1", Message: "A"})
	cb.Add(Message{Type: MessageTypeText, From: "client2", Message: "B"})
	cb.Add(Message{Type: MessageTypeText, From: "client3", Message: "C"})

	messages := cb.GetAll()
	if len(messages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(messages))
	}
	if messages[0].Message != "A" || messages[1].Message != "B" || messages[2].Message != "C" {
		t.Errorf("Unexpected message order or content")
	}
}

func TestCircularBufferEmpty(t *testing.T) {
	cb := server.NewCircularBuffer(3)
	messages := cb.GetAll()
	if len(messages) != 0 {
		t.Errorf("Expected 0 messages in empty buffer, got %d", len(messages))
	}
}
