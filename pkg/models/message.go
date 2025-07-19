package models

type MessageType string 

type Message struct {
	Type MessageType `json:"type"`
	From string `json:"from,omitempty"`
	Message string `json:"message,omitempty"`
}
