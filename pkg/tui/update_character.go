package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/thatstoasty/character-sheet-ui/pkg/database"
)

const (
	selectCharacter State = iota
	selectAttribute
	updatePrompt
)

func updateCharacterAttribute(name string, attribute string, value string) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open(os.Getenv("SHEETS_DATABASE")), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		db.Table("characters").Where("name = ?", name).Update(attribute, value)

		return SwitchStateMsg(selectAttribute)
	}
}

type CharacterMsg database.Character

func getCharacter(name string) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open(os.Getenv("SHEETS_DATABASE")), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		var character database.Character
		db.Table("characters").Where("name = ?", name).Scan(&character)

		return CharacterMsg(character)
	}
}

type UpdateCharacterModel struct {
	State             State
	TextInput         textinput.Model
	List              list.Model
	Character         database.Character
	SelectedCharacter string
	SelectedAttribute string
}

func (m UpdateCharacterModel) Init() tea.Cmd {
	return getCharacterNames
}

func (m UpdateCharacterModel) characterView() string {
	if m.Character.Name == "" {
		return characterStyle.Render("You'll see more here after giving your character a name!\n\n")
	} else {
		// Render the string that represents the character view
		return characterStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
			fmt.Sprintf("%s: The %s", m.Character.Name, m.Character.Race),
			"------------",
			fmt.Sprintf("%d   %d   %d   %d   %d   %d   %d", m.Character.HP, m.Character.Strength, m.Character.Dexterity, m.Character.Constitution, m.Character.Intelligence, m.Character.Wisdom, m.Character.Charisma),
			"HP  STR  DEX  CON  INT  WIS  CHA\n",
			fmt.Sprintf("Class    | %s", m.Character.Class),
			fmt.Sprintf("Feats    | %s", m.Character.Feats),
			fmt.Sprintf("Items    | %s", m.Character.Items),
			fmt.Sprintf("Armor    | %s, %s, %s, %s, %s", m.Character.Helmet, m.Character.Cloak, m.Character.Armor, m.Character.Boots, m.Character.Gloves),
			fmt.Sprintf("Jewelery | %s, %s, %s", m.Character.Jewelery1, m.Character.Jewelery2, m.Character.Jewelery3),
			fmt.Sprintf("Weapons  | %s, %s", m.Character.MainHandWeapon, m.Character.OffHandWeapon),
		))
	}
}

func (m UpdateCharacterModel) promptView() string {
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

func (m UpdateCharacterModel) View() string {
	switch m.State {
	case selectCharacter:
		return m.promptView()
	default:
		return lipgloss.JoinHorizontal(lipgloss.Top, m.promptView(), m.characterView())
	}
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

		case CharacterMsg:
			m.Character = database.Character(msg)
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
				return m, tea.Sequence(setupList("What would you like to update?", &choices), getCharacter(m.SelectedCharacter), switchState(selectAttribute))
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
