package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#aeed43ff")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#95fd4bff")).
			Padding(1, 3).
			Align(lipgloss.Center)
)

type model struct {
	width     int
	height    int
	textInput textinput.Model
	backend   *Backend
	loggedIn  bool
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter your name"
	ti.Focus()
	ti.CharLimit = 16
	ti.Width = 20

	return model{
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.backend != nil {
				m.backend.Close()
			}
			return m, tea.Quit
		case "enter":
			if !m.loggedIn {
				username := m.textInput.Value()
				if username == "" {
					return m, nil
				}
				b, err := NewBackend(username)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.backend = b
				m.backend.StartDiscovery()
				m.loggedIn = true
				// just restarting textinput
				ti := textinput.New()
				ti.Placeholder = "Enter your message"
				ti.Focus()
				ti.CharLimit = 200
				ti.Width = 60
				m.textInput = ti

				return m, nil
			} else {
				msg := strings.TrimSpace(m.textInput.Value())
				if msg == "" {
					return m, nil
				}
				// just restarting textinput
				ti := textinput.New()
				ti.Placeholder = "Enter your message"
				ti.Focus()
				ti.CharLimit = 200
				ti.Width = 60
				m.textInput = ti
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m model) View() string {
	var content string

	if m.err != nil {
		content = fmt.Sprintf("Error: %v\n\nPress 'q' to exit.", m.err)
	} else if !m.loggedIn {
		title := titleStyle.Render("LOCAL P2P CHAT")
		instruction := "Enter your name to start discovering peers:"
		content = fmt.Sprintf("%s\n\n%s\n\n%s", title, instruction, m.textInput.View())
	} else {
		title := titleStyle.Render("LOCAL P2P CHAT")

		content = fmt.Sprintf("%s%s\n\n", title, m.backend.ipHash)
		for _, msg := range m.backend.Messages {
			content += msg.sender + " > " + msg.text
		}

		content += "\n\n" + m.textInput.View()
	}

	box := boxStyle.Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
