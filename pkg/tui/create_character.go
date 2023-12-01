package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	// "log"
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
)

type CreateCharacterModel struct {
	TextInput textinput.Model
	Err       error
}

func (m CreateCharacterModel) Init() tea.Cmd {
	return nil
}

func (m CreateCharacterModel) View() string {
	return fmt.Sprintf(
		"Enter the name of the character to delete:\n\n%s\n\n%s",
		m.TextInput.View(),
		"(esc to quit)",
	) + "\n"
}

func (m CreateCharacterModel) Update(msg tea.Msg) (CreateCharacterModel, tea.Cmd) {
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
	fmt.Println("create")

	var cmd tea.Cmd
	m.TextInput, cmd = m.TextInput.Update(msg)

	return m, cmd
}
