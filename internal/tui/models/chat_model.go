package models

import (
	"fmt"


	. "github.com/sternerr/termtalk/pkg/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/viewport"
)

type ChatModel struct {
	viewport viewport.Model
	messages []Message
	width    int
	height   int
}

func NewChatModel() ChatModel {
	return ChatModel{
		messages: make([]Message, 0, 50),
		viewport: viewport.Model{},
	}
}

func(m ChatModel) Init() tea.Cmd {
	return nil
}

func(m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	switch msg := msg.(type) {

	case Message:
		m.messages = append(m.messages, msg)
		m.updateViewportContent()

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = m.width - 2
    	m.viewport.Height = m.height - 2

		return m, nil
	}

	return m, nil
}

func (m *ChatModel) updateViewportContent() {
	str := ""
	for _, m := range m.messages {
		str += fmt.Sprintf("[%s] %s", m.From, m.Message)
	}

	m.viewport.SetContent(str)
	m.viewport.GotoBottom()
}

func(m ChatModel) View() string {
	viewportStyle := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - 2).
		Border(lipgloss.NormalBorder(), true).
		Padding(0, 1)

	return viewportStyle.Render(m.viewport.View())
}

