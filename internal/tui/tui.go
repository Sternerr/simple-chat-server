package tui

import (
	"log"

	"github.com/sternerr/termtalk/internal/client"
	. "github.com/sternerr/termtalk/pkg/models"
	"github.com/sternerr/termtalk/internal/tui/models"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelTypes int

const (
	ChatModelType ModelTypes = iota
)

type TUI struct {
	Client client.Client
	user User
	messages []Message
	logger *log.Logger
	models map[ModelTypes]tea.Model
}

func NewTUI(logger *log.Logger) TUI {
	return TUI{
		Client: client.NewClient(logger),
		user: User{Username: "A"},
		logger: logger,
		models: map[ModelTypes]tea.Model{
			ChatModelType: models.NewChatModel(),
		},
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

		model, cmd := t.models[ChatModelType].Update(msg)
		t.models[ChatModelType] = model
		return t, tea.Batch(t.ClientListenerCmd(), cmd)
	
	case tea.KeyMsg:
		switch msg.String() {
			case "ctrl+c":
				return t, tea.Quit
		}

	case tea.WindowSizeMsg:
		model, cmd := t.models[ChatModelType].Update(msg)
		t.models[ChatModelType] = model
		return t, cmd
	}

	return t, nil
}

func(t TUI) View() string {
	return t.models[ChatModelType].View()
}

