package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"log"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type RefreshMsg bool
type ListMsg list.Model

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

func setupInitialModel(title string, characterNames *[]string) tea.Cmd {
	return func() tea.Msg {
		items := []list.Item{}

		const defaultWidth = 80

		l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
		l.Title = title
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = titleStyle
		l.Styles.PaginationStyle = paginationStyle
		l.Styles.HelpStyle = helpStyle

		for index, choice := range *characterNames {
			l.InsertItem(index, item(choice))
		}

		return ListMsg(l)
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
		return m, setupInitialModel("Which character do you want to delete?", &names)

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
