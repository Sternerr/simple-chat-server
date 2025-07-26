package tui

import (
	"log"

	"github.com/sternerr/termtalk/internal/client"
	. "github.com/sternerr/termtalk/pkg/models"
	"github.com/sternerr/termtalk/internal/tui/models"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelType int

const (
	DisplayNameModelType ModelType = iota
	ChatModelType
)

type TUI struct {
	Client client.Client
	user User
	messages []Message
	logger *log.Logger
	models map[ModelType]tea.Model
	activeModel ModelType

}

func NewTUI(logger *log.Logger) TUI {
	user := User{Username: "A"}
	return TUI{
		Client: client.NewClient(logger),
		user: user,
		logger: logger,
		models: map[ModelType]tea.Model{
			DisplayNameModelType: models.NewDisplayNameModel(),
			ChatModelType: models.NewChatModel(user),
		},
		activeModel: DisplayNameModelType,
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
	var cmd tea.Cmd
	var model tea.Model

	switch msg {
	case MessageTypeHandshakeDeny:
		(*(t).logger).Println("[info] Handshake Denied")

	case MessageTypeHandshakeAccept:
		(*(t).logger).Println("[info] Handshake Accepted")
		return t, t.ClientListenerCmd()
	}

	switch msg := msg.(type) {
	case client.Connected:
		(*(t).logger).Printf("[info] connected to server: %v\n", msg)
		return t, nil

	case Message:
		(*(t).logger).Printf("[info] Message Recieved: %v\n", msg)

		model, cmd = t.models[t.activeModel].Update(msg)
		t.models[t.activeModel] = model
		return t, tea.Batch(t.ClientListenerCmd(), cmd)

	case models.SendMessageCmd:
		(*(t.logger)).Printf("[info] Sending Message: %s\n", msg.Message)
		t.Client.SendMessage(msg.Message)
		return t, t.ClientListenerCmd()
	
	case models.SendUserCmd:
		t.user = msg.User
		t.activeModel = ChatModelType
		t.Client.SendHandshake(t.user)	
		return t, tea.Batch(t.ClientListenerCmd(), tea.WindowSize())

	case tea.KeyMsg:
		switch msg.String() {
			default:
				model, cmd = t.models[t.activeModel].Update(msg)
				t.models[t.activeModel] = model
			return t, cmd
		}

	case tea.WindowSizeMsg:
		model, cmd := t.models[t.activeModel].Update(msg)
		t.models[t.activeModel] = model
		return t, cmd
	}

	return t, nil
}

func(t TUI) View() string {
	return t.models[t.activeModel].View()
}

