package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	// "github.com/charmbracelet/lipgloss"
	"github.com/thatstoasty/character-sheet-ui/pkg/server"
	"github.com/thatstoasty/character-sheet-ui/pkg/tui"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

func setupDB() tea.Msg {
	server.SetupDB()

	return nil
}

func getCharacterNames() tea.Msg {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	var names []string
	db.Table("characters").Select("name").Scan(&names)
	return CharacterNamesMsg(names)
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
	CharacterNames  []string
}

type Item struct {
	title, desc string
}
type CharacterNamesMsg []string

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

func initialModel() Model {
	items := []list.Item{
		Item{title: "Create Character", desc: "Create your character through interactive prompts!"},
		Item{title: "Delete Character", desc: "Delete an existing character"},
		Item{title: "Update Character", desc: "Update an existing character"},
		Item{title: "Start!", desc: "Starts the web application to generate your interactive character sheet interface!"},
	}

	return Model{
		State: showHome,
		Home: tui.HomeModel{
			List: list.New(items, list.NewDefaultDelegate(), 0, 0),
		},
	}
	// return Model{
	// 	State: showHome,
	// 	Home: tui.HomeModel{
	// 		// Our to-do list is a grocery list
	// 		Choices: []string{"Create Character", "Delete Character", "Update Character", "Start!"},
	// 	},
	// }
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		setupDB,
		getCharacterNames,
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

	case CharacterNamesMsg:
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
			m.DeleteCharacter.Init()
			m.State = 2
		case "Update Character":
			m.Home.Selected = ""
			m.UpdateCharacter.Init()
			m.State = 3
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
