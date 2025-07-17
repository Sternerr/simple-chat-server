package server

import (
	. "github.com/sternerr/termtalk/pkg/models"
)

type CircularBuffer struct {
	F, R int
	Size int
	Buffer []Message
}

func (cb *CircularBuffer) Add(msg Message) {
	cb.R = (cb.R + 1) % cb.Size
	if cb.R == cb.F {
		cb.F = (cb.F + 1) % cb.Size
	}

	cb.Buffer[cb.R] = msg
}

func (cb *CircularBuffer) GetAll() []Message {
	messages := make([]Message, cb.Size)
	if cb.R == -1 {
		return messages
	}
	
	for i := cb.F; i != cb.R; i = (i + 1) % cb.Size {
		messages = append(messages, cb.Buffer[i])
	}

	return messages
}
