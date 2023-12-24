package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

		return SwitchStateMsg(selectAttribute)
	}
}

type UpdateCharacterModel struct {
	State             State
	TextInput         textinput.Model
	List              list.Model
	SelectedCharacter string
	SelectedAttribute string
}

func (m UpdateCharacterModel) Init() tea.Cmd {
	return getCharacterNames
}

func (m UpdateCharacterModel) View() string {
	var prompt string
	if m.State != updatePrompt {
		prompt = fmt.Sprintf("%s\n\n%s\n",
			m.List.View(),
			"(esc to quit, tab to return home)",
		)
	} else {
		prompt = fmt.Sprintf("What would you like to update your %s to?\n\n%s\n\n%s",
			m.SelectedAttribute,
			m.TextInput.View(),
			"(esc to quit, tab to return home)",
		)
	}

	return centeredStyle.Render(prompt)
}

func (m UpdateCharacterModel) Update(msg tea.Msg) (UpdateCharacterModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case SwitchStateMsg:
		m.State = State(msg)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "tab":
			m.State = selectCharacter
			m.SelectedAttribute = ""
			return m, switchParentState(showHome)
		}
	}

	switch m.State {
	case selectCharacter:
		switch msg := msg.(type) {
		case SwitchStateMsg:
			m.State = State(msg)

		case CharacterNamesMsg:
			names := []string(msg)
			return m, setupList("Which character do you want to update?", &names)

		case ListMsg:
			m.List = list.Model(msg)
			return m, nil

		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "ctrl+c":
				return m, tea.Quit

			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.SelectedCharacter = string(i)
				}

				choices := []string{"Race", "HP", "Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma", "Class", "Feats", "Items", "Helmet", "Cloak", "Armor", "Gloves", "Boots", "Jewelery1", "Jewelery2", "Jewelery3", "MainHandWeapon", "OffhandWeapon"}
				return m, tea.Sequence(setupList("What would you like to update?", &choices), switchState(selectAttribute))
			}
		}

	case selectAttribute:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.SelectedAttribute = string(i)
				}

				return m, switchState(updatePrompt)
			}
		}

	case updatePrompt:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {

			case "ctrl+c", "esc":
				return m, tea.Quit

			case "enter":
				m.TextInput.Placeholder = ""
				return m, updateCharacterAttribute(m.SelectedCharacter, m.SelectedAttribute, m.TextInput.Value())
			}
		}

		m.TextInput, cmd = m.TextInput.Update(msg)
		return m, cmd
	}

	m.List, cmd = m.List.Update(msg)
	return m, cmd
}
