package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	initializing State = iota
	selection
	confirmation
)

type CharacterNamesMsg []string

func fetchCharacterNames() tea.Msg {
	db := getDatabaseSession()
	names := getCharacterNames(db)
	return CharacterNamesMsg(names)
}

type RefreshMsg bool

func removeCharacter(name string) tea.Cmd {
	return func() tea.Msg {
		db := getDatabaseSession()
		deleteCharacter(db, name)
		return RefreshMsg(true)
	}
}

type DeleteCharacterModel struct {
	State             State
	List              list.Model
	SelectedCharacter string
}

func (m DeleteCharacterModel) Init() tea.Cmd {
	return fetchCharacterNames
}

func (m DeleteCharacterModel) Update(msg tea.Msg) (DeleteCharacterModel, tea.Cmd) {
	switch msg := msg.(type) {
	case SwitchStateMsg:
		m.State = State(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil
	case RefreshMsg:
		return m, fetchCharacterNames
	case CharacterNamesMsg:
		names := []string(msg)
		return m, setupList("Which character do you want to delete?", &names)
	case ListMsg:
		m.List = list.Model(msg)
		return m, switchState(selection)
	}

	switch m.State {
	case selection:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "ctrl+c":
				return m, tea.Quit

			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.SelectedCharacter = string(i)
				}
				return m, switchState(confirmation)

			case "tab":
				return m, switchParentState(showHome)
			}
		}
	case confirmation:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "ctrl+c":
				return m, tea.Quit

			case "enter":
				character := m.SelectedCharacter
				m.SelectedCharacter = ""
				return m, removeCharacter(character)

			case "tab":
				return m, switchParentState(showHome)
			}
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m DeleteCharacterModel) View() string {
	switch m.State {
	case selection:
		return centeredStyle.Render(
			fmt.Sprintf("%s\n\n%s\n",
				m.List.View(),
				"(esc to quit, tab to return home)",
			),
		)
	case confirmation:
		return centeredStyle.Render(
			fmt.Sprintf("%s %s?\n\n%s\n",
				"Are you sure you want to delete",
				m.SelectedCharacter,
				"(enter to confirm, esc to quit, tab to return home)",
			),
		)
	default:
		return "Loading..."
	}
}
