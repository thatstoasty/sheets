package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Define what menu is being shown by using an state constants
type State int

// Assigns an incrementing value to each of these constants. 0, then 1, then 2, etc...
const (
	showHome State = iota
	showCreateCharacter
	showDeleteCharacter
	showUpdateCharacter
)

func startServer() tea.Cmd {
	return func() tea.Msg {
		startWebServer()

		return nil
	}
}

type SwitchStateMsg State

func switchState(state State) tea.Cmd {
	return func() tea.Msg {
		return SwitchStateMsg(state)
	}
}

type SwitchParentStateMsg State

func switchParentState(state State) tea.Cmd {
	return func() tea.Msg {
		return SwitchParentStateMsg(state)
	}
}

func setupDatabase() tea.Msg {
	seedDatabase()

	return nil
}

type Model struct {
	State           State
	List            list.Model
	CreateCharacter CreateCharacterModel
	DeleteCharacter DeleteCharacterModel
	UpdateCharacter UpdateCharacterModel
	CharacterNames  []string
}

type itemWithDescription struct {
	title, desc string
}

func (i itemWithDescription) Title() string       { return i.title }
func (i itemWithDescription) Description() string { return i.desc }
func (i itemWithDescription) FilterValue() string { return i.title }

func initialModel() Model {
	items := []list.Item{
		itemWithDescription{title: "Create Character", desc: "Create your character through interactive prompts!"},
		itemWithDescription{title: "Delete Character", desc: "Delete an existing character!"},
		itemWithDescription{title: "Update Character", desc: "Update an existing character!"},
		itemWithDescription{title: "Start!", desc: "Starts the webpage to show your character sheet!"},
	}

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Main Menu"

	return Model{
		State: showHome,
		List:  list,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Sequence(
		setupDatabase,
		fetchCharacterNames,
	)
}

func (m Model) View() string {
	switch m.State {
	case showCreateCharacter:
		return m.CreateCharacter.View()
	case showDeleteCharacter:
		return m.DeleteCharacter.View()
	case showUpdateCharacter:
		return m.UpdateCharacter.View()
	default:
		return baseStyle.Render(m.List.View())
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

	case SwitchParentStateMsg:
		m.State = State(msg)
	}

	switch m.State {
	case showHome:
		// Add the list component update message to the list of commands
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)

		switch msg := msg.(type) {
		case SwitchStateMsg:
			m.State = State(msg)

		case SwitchParentStateMsg:
			m.State = State(msg)

		case tea.WindowSizeMsg:
			h, v := centeredStyle.GetFrameSize()
			m.List.SetSize(msg.Width-h, msg.Height-v)

			// Is it a key press?
		case tea.KeyMsg:
			switch msg.String() {

			case "ctrl+c", "esc":
				return m, tea.Quit

			// the selected state for the item that the cursor is pointing at.
			case "enter":
				i, ok := m.List.SelectedItem().(itemWithDescription)
				if ok {
					switch i.Title() {
					case "Create Character":
						return m, tea.Sequence(switchState(showCreateCharacter), m.CreateCharacter.Init())
					case "Delete Character":
						return m, tea.Sequence(switchState(showDeleteCharacter), m.DeleteCharacter.Init())
					case "Update Character":
						return m, tea.Sequence(switchState(showUpdateCharacter), m.UpdateCharacter.Init())
					case "Start!":
						return m, startServer()
					}
				}
			}
		}
		return m, tea.Batch(cmds...)
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

// Start the bubbletea TUI application
func startTUI() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
