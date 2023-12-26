package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"fmt"
	"io"
	"strings"
)

// Lipgloss styles for component rendering
var baseStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("63")).
	PaddingTop(2).
	PaddingBottom(2).
	PaddingLeft(4).
	PaddingRight(4)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	centeredStyle     = baseStyle.Align(lipgloss.Center)
)

// Create a bubbles list with the given string items
func createListWithItems(choices *[]string) list.Model {
	items := []list.Item{}
	for _, choice := range *choices {
		items = append(items, item(choice))
	}

	const listWidth = 80
	const listHeight = 14

	return list.New(items, itemDelegate{}, listWidth, listHeight)
}

type ListMsg list.Model

// Tea CMD: Create a bubbles list with the given string items with styling and a title.
func setupList(title string, characterNames *[]string) tea.Cmd {
	return func() tea.Msg {
		l := createListWithItems(characterNames)

		l.Title = title
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = titleStyle
		l.Styles.PaginationStyle = paginationStyle
		l.Styles.HelpStyle = helpStyle

		return ListMsg(l)
	}
}

// Simple bubbles list item definition with no filtering. List item types must implement the list.Item interface.
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

// Build a text input component with the given prompt and width
func createTextInput() textinput.Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 80
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D70EFF"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E978FF"))

	return ti
}

type TextInputMsg textinput.Model

// Tea CMD: Create a bubbles text input component.
func setupTextInput() tea.Msg {
	return TextInputMsg(createTextInput())
}
