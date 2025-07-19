package models

type MessageType string 

const (
	MessageTypeHandshake MessageType = "handshake"
	MessageTypeHandshakeDeny MessageType = "handshake/deny"
	MessageTypeHandshakeAccept MessageType = "handshake/accept"

	MessageTypeText MessageType = "text"
)

type Message struct {
	Type MessageType `json:"type"`
	From string `json:"from,omitempty"`
	Message string `json:"message,omitempty"`
}
