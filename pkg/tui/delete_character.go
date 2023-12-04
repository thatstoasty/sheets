package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RefreshMsg bool

func deleteCharacter(name string) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		db.Where("name = ?", name).Delete(&Character{})

		return RefreshMsg(true)
	}
}

type DeleteCharacterModel struct {
	Cursor            int // which list choice our cursor is pointing at
	CharacterNames    []string
	SelectedCharacter string
}

func (m DeleteCharacterModel) Init() tea.Cmd {
	return getCharacterNames
}

func (m DeleteCharacterModel) View() string {
	// The header
	s := "Which character would you like to delete?\n\n"

	// Iterate over our choices
	for i, choice := range m.CharacterNames {

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

func (m DeleteCharacterModel) Update(msg tea.Msg) (DeleteCharacterModel, tea.Cmd) {
	switch msg := msg.(type) {

	case RefreshMsg:
		return m, m.Init()

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
			if m.Cursor < len(m.CharacterNames)-1 {
				m.Cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			m.SelectedCharacter = m.CharacterNames[m.Cursor]
			fmt.Println(m.SelectedCharacter)
			m.Cursor = 0
			return m, deleteCharacter(m.SelectedCharacter)
		}
	}

	return m, nil
}
