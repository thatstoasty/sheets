package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thatstoasty/character-sheet-ui/pkg/server"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func startServer() tea.Cmd {
	return func() tea.Msg {
		server.Start()

		return "Server started!"
	}
}

type HomeModel struct {
	List     list.Model
	Selected string // which choice is selected
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m HomeModel) View() string {
	return docStyle.Render(m.List.View())
}

func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit

		// the selected state for the item that the cursor is pointing at.
		case "enter":
			switch m.List.Cursor() {
			case 0:
				m.Selected = "Create Character"
				return m, nil
			case 1:
				m.Selected = "Delete Character"
				return m, nil
			case 2:
				m.Selected = "Update Character"
				return m, nil
			case 3:
				m.Selected = "Start!"
				return m, startServer()
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}
