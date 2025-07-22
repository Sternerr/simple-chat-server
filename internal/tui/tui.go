package tui

import (
	"log"

	. "github.com/sternerr/termtalk/pkg/models"
	tea "github.com/charmbracelet/bubbletea"
)

type TUI struct {
	user User
	logger *log.Logger
}

func NewTUI(logger *log.Logger) TUI {
	return TUI{
		user: User{Username: "A"},
		logger: logger,
	}
}

func(t TUI) Init() tea.Cmd {
	return nil
}

func(t TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
			case "ctrl+c":
				return t, tea.Quit
		}
	}
	return t, nil
}

func(t TUI) View() string {
	return "Chat"
}

