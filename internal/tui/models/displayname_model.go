package models

import (
	"fmt"

	. "github.com/sternerr/termtalk/pkg/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
)

type SendUserCmd struct {
	User User
}

type DisplayNameModel struct {
	input textinput.Model
	user User
	width    int
	height   int
}

func NewDisplayNameModel() DisplayNameModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()

	ti.CharLimit = 32

	return DisplayNameModel{
		input: ti,
	}
}

func(m DisplayNameModel) Init() tea.Cmd {
	return textinput.Blink
}

func(m DisplayNameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	var tiCmd tea.Cmd

	m.input, tiCmd = m.input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "enter":
				user := User{Username: m.input.Value()}
				m.input.Reset()

				return m, func() tea.Msg { return SendUserCmd{User: user}}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = int(float64(msg.Width) * 0.5)
	}

	return m, tiCmd
}

func(m DisplayNameModel) View() string {
	inputStyle := lipgloss.NewStyle().
		Width(m.input.Width).
		Border(lipgloss.NormalBorder(), true)
	
	str := fmt.Sprintf("%s\n%s", "Enter a display name", inputStyle.Render(m.input.View()))
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, str)
}

