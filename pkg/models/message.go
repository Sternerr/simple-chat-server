package models

type Message struct {
	Type string `json:"type"`
	From string `json:"from,omitempty"`
	Message string `json:"message,omitempty"`
}
