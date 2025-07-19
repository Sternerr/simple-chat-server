package server

import . "github.com/sternerr/termtalk/pkg/models"

type CircularBuffer struct {
	F, R int
	Size int
	Count int
	Buffer []Message
}

func(cb *CircularBuffer) Add(msg Message) {
	(*cb).R = ((*cb).R + 1) % cb.Size
	if (*cb).Count == (*cb).Size {
		(*cb).F = ((*cb).F + 1) % (*cb).Size
	} else {
		(*cb).Count++
	}
	cb.Buffer[(*cb).R] = msg
}

func (cb *CircularBuffer) GetAll() []Message {
	messages := make([]Message, 0, (*cb).Count)
	for i, count := (*cb).F, 0; count < (*cb).Count; count++ {
		messages = append(messages, (*cb).Buffer[i])
		i = (i + 1) % (*cb).Size
	}
	return messages
}

func NewCircularBuffer(size int) CircularBuffer {
	return CircularBuffer{
		F:      0,
		R:      -1,
		Size:   size,
		Count:  0,
		Buffer: make([]Message, size),
	}
}
