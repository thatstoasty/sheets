package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thatstoasty/character-sheet-ui/pkg/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func startServer() tea.Cmd {
	return func() tea.Msg {
		server.Start()

		return "Server started!"
	}
}

func switchState(state State) tea.Cmd {
	return func() tea.Msg {
		return SwitchStateMsg(state)
	}
}

type SwitchStateMsg State

// Define what menu is being shown by using an state constants
type State int

func setupDB() tea.Msg {
	server.SetupDB()

	return nil
}

func getCharacterNames() tea.Msg {
	fmt.Println("getting character names")

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

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// Assigns an incrementing value to each of these constants. 0, then 1, then 2, etc...
const (
	showHome State = iota
	showCreateCharacter
	showDeleteCharacter
	showUpdateCharacter
)

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
type CharacterNamesMsg []string

func (i itemWithDescription) Title() string       { return i.title }
func (i itemWithDescription) Description() string { return i.desc }
func (i itemWithDescription) FilterValue() string { return i.title }

func initialModel() Model {
	items := []list.Item{
		itemWithDescription{title: "Create Character", desc: "Create your character through interactive prompts!"},
		itemWithDescription{title: "Delete Character", desc: "Delete an existing character"},
		itemWithDescription{title: "Update Character", desc: "Update an existing character"},
		itemWithDescription{title: "Start!", desc: "Starts the web application to generate your interactive character sheet interface!"},
	}

	return Model{
		State: showHome,
		List:  list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		setupDB,
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
		return docStyle.Render(m.List.View())
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.State {
	case showHome:
		// Add the list component update message to the list of commands
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)

		switch msg := msg.(type) {
		case SetupTextInputMsg:
			switch msg.ModelName {
			case "CreateCharacter":
				m.CreateCharacter.TextInput = msg.TextInput
			case "UpdateCharacter":
				m.UpdateCharacter.TextInput = msg.TextInput
			}

		case CreateListMsg:
			m.CreateCharacter.List = list.Model(msg)

		case ListMsg:
			m.DeleteCharacter.List = list.Model(msg)

		case RacesMsg:
			m.CreateCharacter.Races = []string(msg)

		case SwitchStateMsg:
			m.State = State(msg)

		case CharacterNamesMsg:
			m.UpdateCharacter.AttributeChoices = []string{"Race", "HP", "Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma", "Class", "Feats", "Items", "Helmet", "Cloak", "Armor", "Gloves", "Boots", "Jewelery1", "Jewelery2", "Jewelery3", "MainHandWeapon", "OffhandWeapon"}
			m.UpdateCharacter.CharacterNames = []string(msg)
			m.DeleteCharacter.CharacterNames = []string(msg)

		case tea.WindowSizeMsg:
			h, v := docStyle.GetFrameSize()
			m.List.SetSize(msg.Width-h, msg.Height-v)

			// Is it a key press?
		case tea.KeyMsg:
			switch msg.String() {

			case "ctrl+c", "esc":
				return m, tea.Quit

			// the selected state for the item that the cursor is pointing at.
			case "enter":
				switch m.List.Cursor() {
				case 0:
					cmds = append(cmds, m.CreateCharacter.Init(), switchState(showCreateCharacter))
					return m, tea.Sequence(cmds...)
				case 1:
					cmds = append(cmds, m.DeleteCharacter.Init(), switchState(showDeleteCharacter))
					return m, tea.Sequence(cmds...)
				case 2:
					cmds = append(cmds, m.UpdateCharacter.Init(), switchState(showUpdateCharacter))
					return m, tea.Sequence(cmds...)
				case 3:
					return m, startServer()
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

func Run() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
