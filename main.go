package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thatstoasty/character-sheet-ui/pkg/server"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB() tea.Msg {
	server.SetupDB()

	return nil
}

func startServer() tea.Cmd {
	return func() tea.Msg {
		server.Start()

		return "Server started!"
	}
}

func deleteCharacter(name string) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		db.Table("characters").Where("name = ?", name)

		return logMsg(fmt.Sprintf("%s has been deleted!", name))
	}
}

type Model struct {
	Choices         []string         // items on the to-do list
	Cursor          int              // which to-do list item our cursor is pointing at
	Selected        map[int]struct{} // which to-do items are selected
	StartServer     bool
	DeleteCharacter bool
	TextInput       textinput.Model
	Err             error
}

type logMsg string

func initialModel() Model {
	return Model{
		// Our to-do list is a grocery list
		Choices: []string{"Create Character", "Delete Character", "Start!"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		Selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return setupDB
}

func (m Model) View() string {
	// The header
	s := "What would you like to do?\n\n"

	if m.DeleteCharacter {
		return fmt.Sprintf(
			"Enter the name of the character to delete:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	} else {
		// Iterate over our choices
		for i, choice := range m.Choices {

			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.Cursor == i {
				cursor = ">" // cursor!
			}

			// Is this choice selected?
			checked := " " // not selected
			if _, ok := m.Selected[i]; ok {
				checked = "x" // selected!
			}

			// Render the row
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)

		}
	}

	// The footer
	s += "\nPress esc to quit.\n"

	// Send the UI for rendering
	return s
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case logMsg:
		m.DeleteCharacter = false
		m.TextInput = textinput.New()

	// Is it a key press?
	case tea.KeyMsg:
		if m.DeleteCharacter {
			switch msg.String() {

			// These keys should exit the program.
			case "ctrl+c", "esc":
				return m, tea.Quit

			case "enter":
				return m, deleteCharacter(m.TextInput.Value())
			}

			var cmd tea.Cmd
			m.TextInput, cmd = m.TextInput.Update(msg)

			return m, cmd
		} else {
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
				case 0:
					return m, nil
				case 1:
					m.DeleteCharacter = true
					ti := textinput.New()
					ti.Placeholder = "Mikhail"
					ti.Focus()
					ti.CharLimit = 156
					ti.Width = 20
					m.TextInput = ti
					return m, nil
				case 2:
					return m, startServer()
				}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
