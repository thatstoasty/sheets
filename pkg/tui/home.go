package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thatstoasty/character-sheet-ui/pkg/server"
)

func startServer() tea.Cmd {
	return func() tea.Msg {
		server.Start()

		return "Server started!"
	}
}

type HomeModel struct {
	Choices  []string // choices on the list
	Cursor   int      // which list choice our cursor is pointing at
	Selected string   // which choice is selected
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m HomeModel) View() string {
	// The header
	s := "What would you like to do?\n\n"

	// Iterate over our choices
	for i, choice := range m.Choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.Cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	// The footer
	s += "\nPress esc to quit.\n"

	// Send the UI for rendering
	return s
}

func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			switch m.Cursor {
			// case 0:
			// 	m.Selected = "Create Character"
			// 	return m, nil
			case 0:
				m.Selected = "Delete Character"
				return m, nil
			case 1:
				m.Selected = "Start!"
				return m, startServer()
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
