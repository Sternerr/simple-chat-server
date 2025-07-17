package protocol

import (
	"bytes"
	"errors"
	"fmt"
)

func DecodeHandshake(msg []byte) (string, error) {
	header, content, ok := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !ok {
		return "", errors.New("did not find delimiter in handshake")
	}

	trimmedHeader := bytes.TrimSpace(header)
	if !bytes.Equal(trimmedHeader, []byte("Handshake")) {
		return "", errors.New("invalid handshake format")
	}
	
	fmt.Println(content)
	return string(content), nil
}
