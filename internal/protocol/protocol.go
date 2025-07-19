package protocol

import (
	"fmt"
	"bytes"
	"strings"
	"errors"
	"encoding/json"

	. "github.com/sternerr/termtalk/pkg/models"
)

func EncodeMessage(msg Message) ([]byte, error){
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("could not encode msg: %s", err.Error())
	}

	return data, nil
}

func DecodeMessage(msg []byte) (Message, error) {
	msg = bytes.TrimSpace(msg)
	if len(msg) <= 0 {
		return Message{}, errors.New("message is empty")		
	}
	
	var message Message
	err := json.Unmarshal(msg, &message)
	if err != nil {
		return Message{}, errors.New("invalid json")
	}
	
	switch message.Type {
	case "handshake":
		if message.From == "" {
			return Message{}, errors.New("missing 'from' in handshake")
		}
	case "chat":
		if message.From == "" || message.Message == "" {
			return Message{}, errors.New("missing 'from' or 'message' in chat")
		}
	default:
		return Message{}, errors.New("unkown message type")
	}
	
	return message, nil
}

func IsValidHandshake(msg Message) bool {
	from := strings.TrimSpace(msg.From)
	if msg.Type == "" || from == "" {
		return false
	}
	
	if msg.Type != "handshake" {
		return false
	}

	return true
}
