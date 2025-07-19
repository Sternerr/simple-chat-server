package models

type MessageType string 

const (
	MessageTypeHandshake MessageType = "handshake"
	MessageTypeHandshakeDeny MessageType = "handshake/deny"
)

type Message struct {
	Type MessageType `json:"type"`
	From string `json:"from,omitempty"`
	Message string `json:"message,omitempty"`
}
