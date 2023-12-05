package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CharacterMsg Character
type SwitchUpdateStateMsg State

const (
	selectCharacter State = iota
	selectAttribute
	updatePrompt
)

func updateCharacterAttribute(name string, attribute string, value string) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		db.Table("characters").Where("name = ?", name).Update(attribute, value)

		return SwitchUpdateStateMsg(0)
	}
}

func getCharacter(name string) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		var character Character
		db.First(&character, "name = ?", name)

		return CharacterMsg(character)
	}
}

type UpdateCharacterModel struct {
	State             State
	TextInput         textinput.Model
	Err               error
	Cursor            int // which list choice our cursor is pointing at
	CharacterNames    []string
	SelectedCharacter string
	Character         Character
	AttributeChoices  []string
	SelectedAttribute string
}

func (m UpdateCharacterModel) Init() tea.Cmd {
	return nil
}

func (m UpdateCharacterModel) View() string {
	switch m.State {
	case selectCharacter:
		// The header
		s := "Which character would you like to update?\n\n"

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
	case selectAttribute:
		// The header
		s := "Which would you like to update?\n\n"

		// Iterate over our choices
		for i, choice := range m.AttributeChoices {

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
	case updatePrompt:
		return fmt.Sprintf(
			"What would you like to update your %s to?:\n\n%s\n\n%s",
			m.SelectedAttribute,
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	}

	return "Oops!"
}

func (m UpdateCharacterModel) Update(msg tea.Msg) (UpdateCharacterModel, tea.Cmd) {
	switch m.State {
	case selectCharacter:

		switch msg := msg.(type) {
		case CharacterMsg:
			m.Character = Character(msg)
			m.State = 1
			return m, nil

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
				m.Cursor = 0
				return m, getCharacter(m.SelectedCharacter)
			}
		}

	case selectAttribute:

		switch msg := msg.(type) {

		case SwitchUpdateStateMsg:
			m.State = State(msg)

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
				if m.Cursor < len(m.AttributeChoices)-1 {
					m.Cursor++
				}

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			case "enter":
				m.SelectedAttribute = m.AttributeChoices[m.Cursor]
				m.Cursor = 0
				m.TextInput.Placeholder = m.SelectedAttribute
				m.State = 2 //Switch to update prompt
				return m, nil
			}
		}

	case updatePrompt:
		switch msg := msg.(type) {

		case SwitchUpdateStateMsg:
			m.State = State(msg)

		// Is it a key press?
		case tea.KeyMsg:
			// Cool, what was the actual key pressed?
			switch msg.String() {

			// These keys should exit the program.
			case "ctrl+c", "esc":
				return m, tea.Quit

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			case "enter", " ":
				m.TextInput.Placeholder = ""
				return m, updateCharacterAttribute(m.SelectedCharacter, m.SelectedAttribute, m.TextInput.Value())
			}
		}

		// TextInput component Update -> DeleteCharacter component Update -> Parent Update. Bubbles up
		var cmd tea.Cmd
		m.TextInput, cmd = m.TextInput.Update(msg)

		return m, cmd
	}

	return m, nil
}
