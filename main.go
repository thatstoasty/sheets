package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thatstoasty/character-sheet-ui/pkg/server"
	"github.com/thatstoasty/character-sheet-ui/pkg/tui"
	"os"
)

func setupDB() tea.Msg {
	server.SetupDB()

	return nil
}

// Define what menu is being shown by using an state constants
type State int

// Assigns an incrementing value to each of these constants. 0, then 1, then 2, etc...
const (
	showHome State = iota
	// showCreateCharacter
	showDeleteCharacter
)

type Model struct {
	State State
	Home  tui.HomeModel
	// CreateCharacter tui.CreateCharacterModel
	DeleteCharacter tui.DeleteCharacterModel
}

func initialModel() Model {
	return Model{
		State: showHome,
		Home: tui.HomeModel{
			// Our to-do list is a grocery list
			Choices: []string{"Delete Character", "Start!"},
		},
	}
}

func (m Model) Init() tea.Cmd {
	return setupDB
}

func (m Model) View() string {
	switch m.State {
	case showHome:
		return m.Home.View()
	// case showCreateCharacter:
	// 	return m.CreateCharacter.View()
	case showDeleteCharacter:
		return m.DeleteCharacter.View()
	default:
		return m.Home.View()
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.State {
	case showHome:
		m.Home, cmd = m.Home.Update(msg)
		switch m.Home.Selected {
		// If Create Character is selected (1st position in the list), then switch to state 1 (create character menu).
		// case "Create Character":
		// 	m.State = 2
		// If Delete Character is selected (2nd position in the list), then switch to state 2 (delete character menu).
		case "Delete Character":
			m.State = 1
		}
		return m, cmd
	// case showCreateCharacter:
	// 	m.CreateCharacter, cmd = m.CreateCharacter.Update(msg)
	// 	return m, cmd
	case showDeleteCharacter:
		m.DeleteCharacter, cmd = m.DeleteCharacter.Update(msg)
		return m, cmd
	default:
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
