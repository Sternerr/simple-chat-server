package models

import (
	"fmt"


	. "github.com/sternerr/termtalk/pkg/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbles/textarea"
)

type SendMessageCmd struct {
	Message Message
}

type ChatModel struct {
	viewport viewport.Model
	textarea textarea.Model
	messages []Message
	user User
	width    int
	height   int
}

func NewChatModel(user User) ChatModel {
	ta := textarea.New()
	ta.Placeholder = ""
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return ChatModel{
		messages: make([]Message, 0, 50),
		viewport: viewport.Model{},
		textarea: ta,
		user: user,
	}
}

func(m ChatModel) Init() tea.Cmd {
	return textarea.Blink
}

func(m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	var taCmd tea.Cmd

	m.textarea, taCmd = m.textarea.Update(msg)

	switch msg := msg.(type) {
	case Message:
		m.messages = append(m.messages, msg)
		m.updateViewportContent()

	case tea.KeyMsg:
		switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "enter":
				message := Message{Type: "text", From: m.user.Username, Message: m.textarea.Value() + "\n"}
				m.messages = append(m.messages, message)
				m.updateViewportContent()
				m.textarea.Reset()

				return m, func() tea.Msg { return SendMessageCmd{Message: message}}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = m.width - 2
		m.textarea.SetWidth(msg.Width)
		m.textarea.SetHeight(3)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height("\n\n\n")
	}

	return m, taCmd
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
		Border(lipgloss.NormalBorder(), true)

	textareaStyle := lipgloss.NewStyle().
		Width(m.width - 2).
		Border(lipgloss.NormalBorder(), true)

	return fmt.Sprintf(
		"%s%s%s",
		viewportStyle.Render(m.viewport.View()),
		"\n",
		textareaStyle.Render(m.textarea.View()),
	)
}

