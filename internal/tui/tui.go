package tui

import (
	"fmt"
	"log"

	"github.com/sternerr/termtalk/internal/client"
	. "github.com/sternerr/termtalk/pkg/models"
	tea "github.com/charmbracelet/bubbletea"
)

type TUI struct {
	Client client.Client
	user User
	messages []Message
	logger *log.Logger
}

func NewTUI(logger *log.Logger) TUI {
	return TUI{
		Client: client.NewClient(logger),
		user: User{Username: "A"},
		messages: make([]Message, 0, 50),
		logger: logger,
	}
}

func(t TUI) ClientListenerCmd() tea.Cmd {
	return func() tea.Msg {
		return <-t.Client.MsgChan
	}
}

func(t TUI) Init() tea.Cmd {
	return t.ClientListenerCmd()
}

func(t TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	(*(t).logger).Println(msg)

	switch msg {
	case MessageTypeHandshakeDeny:
		(*(t).logger).Println("[info] Handshake Denied")

	case MessageTypeHandshakeAccept:
		(*(t).logger).Println("[info] Handshake Accepted")
		return t, t.ClientListenerCmd()
	}

	switch msg := msg.(type) {
	case client.Connected:
		t.Client.SendHandshake(t.user)	
		return t, t.ClientListenerCmd()

	case Message:
		(*(t).logger).Printf("[info] Message Recieved: %s\n", msg)
		t.messages = append(t.messages, msg)
		return t, t.ClientListenerCmd()
	
	case tea.KeyMsg:
		switch msg.String() {
			case "ctrl+c":
				return t, tea.Quit
		}
	}

	return t, nil
}

func(t TUI) View() string {
	str := ""
	for _, m := range t.messages {
		str += fmt.Sprintf("[%s] %s", m.From, m.Message)
	}

	return str
}

