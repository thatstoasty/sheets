package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type CharacterNamesMsg []string

func getCharacterNames() tea.Msg {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	var names []string
	db.Table("characters").Select("name").Scan(&names)
	return CharacterNamesMsg(names)
}

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
	Initalizing bool
	List        list.Model
}

func (m DeleteCharacterModel) Init() tea.Cmd {
	return getCharacterNames
}

func (m DeleteCharacterModel) Update(msg tea.Msg) (DeleteCharacterModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ListMsg:
		m.List = list.Model(msg)
		m.Initalizing = false

	case RefreshMsg:
		return m, getCharacterNames

	case CharacterNamesMsg:
		names := []string(msg)
		return m, setupList("Which character do you want to delete?", &names)

	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			var cmd tea.Cmd
			i, ok := m.List.SelectedItem().(item)
			if ok {
				cmd = deleteCharacter(string(i))
			}

			return m, cmd

		case "tab":
			return m, switchParentState(showHome)
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m DeleteCharacterModel) View() string {
	if m.Initalizing {
		return "Loading..."
	} else {
		return fmt.Sprintf("%s\n\n%s\n",
			m.List.View(),
			"(esc to quit, tab to return home)",
		)
	}
}
