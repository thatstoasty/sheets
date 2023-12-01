package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

type DeleteCharacterModel struct {
	TextInput textinput.Model
	Err       error
}

type logMsg string

func (m DeleteCharacterModel) Init() tea.Cmd {
	return nil
}

func (m DeleteCharacterModel) View() string {
	return fmt.Sprintf(
		"Enter the name of the character to delete:\n\n%s\n\n%s",
		m.TextInput.View(),
		"(esc to quit)",
	) + "\n"
}

func (m DeleteCharacterModel) Update(msg tea.Msg) (DeleteCharacterModel, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			return m, deleteCharacter(m.TextInput.Value())
		}
	}

	// TextInput component Update -> DeleteCharacter component Update -> Parent Update. Bubbles up
	var cmd tea.Cmd
	m.TextInput, cmd = m.TextInput.Update(msg)

	return m, cmd
}
