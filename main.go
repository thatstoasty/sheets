package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thatstoasty/character-sheet-ui/pkg/server"
	"github.com/thatstoasty/character-sheet-ui/pkg/tui"
	"os"
)

func setupDB() tea.Msg {
	server.SetupDB()

	return nil
}

func BuildTextInput() textinput.Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return ti
}

// Assigns an incrementing value to each of these constants. 0, then 1, then 2, etc...
const (
	showHome tui.State = iota
	showCreateCharacter
	showDeleteCharacter
	showUpdateCharacter
)

type Model struct {
	State           tui.State
	Home            tui.HomeModel
	CreateCharacter tui.CreateCharacterModel
	DeleteCharacter tui.DeleteCharacterModel
	UpdateCharacter tui.UpdateCharacterModel
}

func initialModel() Model {
	return Model{
		State: showHome,
		Home: tui.HomeModel{
			// Our to-do list is a grocery list
			Choices: []string{"Create Character", "Delete Character", "Update Character", "Start!"},
		},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		setupDB,
		m.UpdateCharacter.Init(),
	)
}

func (m Model) View() string {
	switch m.State {
	case showHome:
		return m.Home.View()
	case showCreateCharacter:
		return m.CreateCharacter.View()
	case showDeleteCharacter:
		return m.DeleteCharacter.View()
	case showUpdateCharacter:
		return m.UpdateCharacter.View()
	default:
		return m.Home.View()
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tui.SwitchStateMsg:
		m.State = tui.State(msg)

	case tui.CharacterNamesMsg:
		m.UpdateCharacter.AttributeChoices = []string{"Race", "HP", "Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma", "Class", "Feats", "Items", "Helmet", "Cloak", "Armor", "Gloves", "Boots", "Jewelery1", "Jewelery2", "Jewelery3", "MainHandWeapon", "OffhandWeapon"}
		m.UpdateCharacter.CharacterNames = []string(msg)
		m.DeleteCharacter.CharacterNames = []string(msg)
	}

	switch m.State {
	case showHome:
		m.Home, cmd = m.Home.Update(msg)
		switch m.Home.Selected {
		// If Create Character is selected (1st position in the list), then switch to state 1 (create character menu).
		case "Create Character":
			m.Home.Selected = ""
			m.State = 1
			m.CreateCharacter.TextInput = BuildTextInput()
		// If Delete Character is selected (2nd position in the list), then switch to state 2 (delete character menu).
		case "Delete Character":
			m.Home.Selected = ""
			m.State = 2
			m.DeleteCharacter.Init()
		case "Update Character":
			m.Home.Selected = ""
			m.State = 3
			m.UpdateCharacter.Init()
			m.UpdateCharacter.TextInput = BuildTextInput()
		}
		return m, cmd
	case showCreateCharacter:
		m.CreateCharacter, cmd = m.CreateCharacter.Update(msg)
		return m, cmd
	case showDeleteCharacter:
		m.DeleteCharacter, cmd = m.DeleteCharacter.Update(msg)
		return m, cmd
	case showUpdateCharacter:
		m.UpdateCharacter, cmd = m.UpdateCharacter.Update(msg)
		return m, cmd
	default:
		m.State = 0
		return m, nil
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
