package tui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/slices"

	// "log"
	"github.com/thatstoasty/character-sheet-ui/pkg/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Name            = "Mikhail"
// Race            = "Dwarf"
// HP	  			= 30
// Strength 		= 16
// Dexterity 		= 10
// Constitution     = 14
// Intelligence 	= 8
// Wisdom 			= 8
// Charisma 		= 17
// Feats           = "War Caster,Polearm Master,Tavern Brawler,Tough,Spell Sniper,Great Weapon Master,Sharpshooter"
// Items           = "Potion of Healing,Shortbow"
// Class           = "Paladin|10|Oath of Ancients,Monk|10|Way of the Shadow"
// Helmet          = "Leather Helmet"
// Cloak           = "Leather Cloak"
// Armor           = "Scale Mail"
// Jewelery1       = "Bronze Ring"
// Jewelery2       = "Bronze Ring"
// Jewelery3       = "Bronze Ring"
// Boots           = "Leather Boots"
// Gloves          = "Leather Gloves"
// MainHandWeapon  = "Rapier"
// OffHandWeapon   = ""

type FeatureWithChoices struct {
	Name    string
	Choices []string
}

func (i FeatureWithChoices) FilterValue() string { return "" }

type ChoicesMsg []FeatureWithChoices

func getChoices(class string) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		classes := strings.Split(class, ",")

		var records []ClassFeature
		var features []FeatureWithChoices
		for _, class := range classes {
			classConfig := strings.Split(class, "|")
			_ = db.Table("class_features").Where("class = ? AND level <= ? AND sub_class in ? AND type = 'Choice'", classConfig[0], classConfig[1], []string{"Base", classConfig[2]}).Find(&records)
			for _, record := range records {
				features = append(features, FeatureWithChoices{record.Name, strings.Split(record.Options, "|")})
			}
		}

		return ChoicesMsg(features)
	}
}

type TableData struct {
	Category string
	Records  []string
}

type TableDataMsg TableData

func getRaces() tea.Msg {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	var races []string
	db.Table("races").Select("name").Scan(&races)
	return TableDataMsg(TableData{"Race", races})
}

func getWeapons() tea.Msg {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	var weapons []string
	db.Table("items").Select("name").Where("category = 'Weapon'").Scan(&weapons)
	return TableDataMsg(TableData{"Weapon", weapons})
}

func getLightWeapons() tea.Msg {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	var weapons []string
	db.Table("items").Select("name").Where("category = 'Weapon' and properties like '%Light%'").Scan(&weapons)
	return TableDataMsg(TableData{"LightWeapon", weapons})
}

func getGear(category string) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		var weapons []string
		db.Table("items").Select("name").Where("type = ?", category).Scan(&weapons)
		return TableDataMsg(TableData{category, weapons})
	}
}

type SubmitCharacterMsg bool

func submitChoices(character string, choices []FeatureWithChoices) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		for _, feature := range choices {
			_ = db.Save(&server.FeatureChoices{Character: character, Feature: feature.Name, Choice: feature.Choices[0]})
		}

		return SubmitCharacterMsg(true)
	}
}

func populateCharacterNames(list list.Model, characterNames []string) tea.Cmd {
	return func() tea.Msg {
		for index, choice := range characterNames {
			list.InsertItem(index, item(choice))
		}

		return ListMsg(list)
	}
}

// Assigns an incrementing value to each of these constants. 0, then 1, then 2, etc...
const (
	promptName State = iota
	promptRace
	promptClass
	promptStats
	promptFeats
	promptItems
	promptHelmet
	promptCloak
	promptArmor
	promptJewelery1
	promptJewelery2
	promptJewelery3
	promptBoots
	promptGloves
	promptMainHandWeapon
	promptOffhandWeapon
	promptSelection
	characterCreated
)

type CreateCharacterModel struct {
	State         State
	TextInput     textinput.Model
	Err           error
	Character     Character
	List          list.Model
	Choices       []FeatureWithChoices
	ChoiceIndex   int
	ActiveFeature string
	ActiveChoices []string
	Selections    []FeatureWithChoices
	Races         []string
	Weapons       []string
	LightWeapons  []string
	Armor         []string
	Cloaks        []string
	Helmets       []string
	Boots         []string
	Gloves        []string
	Jewelery      []string
}

type CreateListMsg list.Model

func setupInitialCreateList(title string, elements *[]string) tea.Cmd {
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

		for index, element := range *elements {
			l.InsertItem(index, item(element))
		}

		return CreateListMsg(l)
	}
}

func submitCharacter(hero Character) tea.Cmd {
	return func() tea.Msg {
		SubmitCharacter(hero)
		return SwitchStateMsg(showHome)
	}
}

type TextInputMsg textinput.Model

func setupTextInput() tea.Msg {
	return TextInputMsg(BuildTextInput())
}

func (m CreateCharacterModel) Init() tea.Cmd {
	return tea.Batch(
		getRaces,
		setupTextInput,
	)

}

func (m CreateCharacterModel) View() string {
	if slices.Contains([]State{promptName, promptClass, promptStats, promptFeats, promptItems, promptItems, characterCreated}, m.State) {
		switch m.State {
		case promptName:
			return fmt.Sprintf(
				"Enter the name of the character:\n\n%s\n\n%s",
				m.TextInput.View(),
				"(esc to quit)",
			) + "\n"
		case promptClass:
			return fmt.Sprintf(
				"Enter the classes of the character (| and , delimited like this Paladin|10|Oath of Ancients,Monk|10|Way of the Shadow):\n\n%s\n\n%s",
				m.TextInput.View(),
				"(esc to quit)",
			) + "\n"
		case promptRace:
			return m.List.View()
		case promptStats:
			return fmt.Sprintf(
				"Enter the stats of the character separated by commas (HP, Strength, Dexterity, Constitution, Intelligence, Wisdom, Charisma):\n\n%s\n\n%s",
				m.TextInput.View(),
				"(esc to quit)",
			) + "\n"
		case promptFeats:
			return fmt.Sprintf(
				"Enter the feats your character has separated by commas:\n\n%s\n\n%s",
				m.TextInput.View(),
				"(esc to quit)",
			) + "\n"
		case promptItems:
			return fmt.Sprintf(
				"Enter the items your character has separated by commas:\n\n%s\n\n%s",
				m.TextInput.View(),
				"(esc to quit)",
			) + "\n"
		case characterCreated:
			return fmt.Sprintf(
				"%s has been created!\n\n%s\n\n%s",
				m.Character.Name,
				"(tab to return home)",
				"(esc to quit)",
			) + "\n"
		}
	} else {
		return m.List.View()
	}

	return "Oops!"
}

func (m CreateCharacterModel) Update(msg tea.Msg) (CreateCharacterModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case CreateListMsg:
		m.List = list.Model(msg)

	case TextInputMsg:
		m.TextInput = textinput.Model(msg)

	case SwitchStateMsg:
		m.TextInput.Reset()
		m.State = State(msg)

	case TableDataMsg:
		tableData := TableData(msg)
		switch tableData.Category {
		case "Weapon":
			m.Weapons = tableData.Records
		case "LightWeapon":
			m.LightWeapons = tableData.Records
		case "Armor":
			m.Armor = tableData.Records
		case "Helmet":
			m.Helmets = tableData.Records
		case "Cloak":
			m.Cloaks = tableData.Records
		case "Boots":
			m.Boots = tableData.Records
		case "Gloves":
			m.Gloves = tableData.Records
		case "Jewelery":
			m.Jewelery = tableData.Records
		case "Race":
			m.Races = tableData.Records
		}
	}

	switch m.State {
	case promptName:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Name = m.TextInput.Value()
				return m, tea.Sequence(setupInitialCreateList("Choose your race", &m.Races), switchState(promptRace))
			}
		}
	case promptRace:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Character.Race = string(i)
				}

				return m, tea.Batch(
					getGear("Cloak"),
					getGear("Helmet"),
					getGear("Armor"),
					getGear("Jewelery"),
					getGear("Boots"),
					getGear("Gloves"),
					getWeapons,
					getLightWeapons,
					switchState(promptClass),
				)
			}
		}
	case promptClass:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Class = m.TextInput.Value()
				return m, switchState(promptStats)
			}
		}
	case promptStats:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				stats := strings.Split(m.TextInput.Value(), ",")
				hp, err := strconv.Atoi(stats[0])
				if err != nil {
					log.Fatal("Failed to convert HP to an integer.")
				}

				strength, err := strconv.Atoi(strings.Trim(stats[1], " "))
				if err != nil {
					log.Fatal("Failed to convert Strength to an integer.")
				}

				dexterity, err := strconv.Atoi(strings.Trim(stats[2], " "))
				if err != nil {
					log.Fatal("Failed to convert Dexterity to an integer.")
				}

				constitution, err := strconv.Atoi(strings.Trim(stats[3], " "))
				if err != nil {
					log.Fatal("Failed to convert Constitution to an integer.")
				}

				intelligence, err := strconv.Atoi(strings.Trim(stats[4], " "))
				if err != nil {
					log.Fatal("Failed to convert Intelligence to an integer.")
				}

				wisdom, err := strconv.Atoi(strings.Trim(stats[5], " "))
				if err != nil {
					log.Fatal("Failed to convert Wisdom to an integer.")
				}

				charisma, err := strconv.Atoi(strings.Trim(stats[6], " "))
				if err != nil {
					log.Fatal("Failed to convert Charisma to an integer.")
				}

				m.Character.HP = uint8(hp)
				m.Character.Strength = uint8(strength)
				m.Character.Dexterity = uint8(dexterity)
				m.Character.Constitution = uint8(constitution)
				m.Character.Intelligence = uint8(intelligence)
				m.Character.Wisdom = uint8(wisdom)
				m.Character.Charisma = uint8(charisma)
				return m, switchState(promptFeats)
			}
		}
	case promptFeats:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Feats = m.TextInput.Value()
				return m, switchState(promptItems)
			}
		}
	case promptItems:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Items = m.TextInput.Value()
				return m, tea.Sequence(setupInitialCreateList("What helmet does your character have equipped?", &m.Helmets), switchState(promptHelmet))
			}
		}
	case promptHelmet:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Helmet = m.TextInput.Value()
				return m, tea.Sequence(setupInitialCreateList("What cloak does your character have equipped?", &m.Cloaks), switchState(promptCloak))
			}
		}
	case promptCloak:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Character.Cloak = string(i)
				}
				return m, tea.Sequence(setupInitialCreateList("What armor does your character have equipped?", &m.Armor), switchState(promptArmor))
			}
		}
	case promptArmor:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Character.Armor = string(i)
				}
				return m, tea.Sequence(setupInitialCreateList("What's the first piece jewelery your character has equipped?", &m.Jewelery), switchState(promptJewelery1))
			}
		}
	case promptJewelery1:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Character.Jewelery1 = string(i)
				}
				return m, tea.Sequence(setupInitialCreateList("What's the second piece jewelery your character has equipped?", &m.Jewelery), switchState(promptJewelery2))
			}
		}
	case promptJewelery2:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Character.Jewelery2 = string(i)
				}
				return m, tea.Sequence(setupInitialCreateList("What's the third piece jewelery your character has equipped?", &m.Jewelery), switchState(promptJewelery3))
			}
		}
	case promptJewelery3:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Character.Jewelery3 = string(i)
				}
				return m, tea.Sequence(setupInitialCreateList("What boots does your character have equipped?", &m.Boots), switchState(promptBoots))
			}
		}
	case promptBoots:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Character.Boots = string(i)
				}
				return m, tea.Sequence(setupInitialCreateList("What gloves does your character have equipped?", &m.Gloves), switchState(promptGloves))
			}
		}
	case promptGloves:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Character.Gloves = string(i)
				}
				return m, tea.Sequence(setupInitialCreateList("What main hand weapon does your character have equipped?", &m.Weapons), switchState(promptMainHandWeapon))
			}
		}
	case promptMainHandWeapon:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Character.MainHandWeapon = string(i)
				}
				return m, tea.Sequence(setupInitialCreateList("What off hand weapon does your character have equipped?", &m.LightWeapons), switchState(promptOffhandWeapon))
			}
		}
	case promptOffhandWeapon:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Character.OffHandWeapon = string(i)
				}
				return m, getChoices(m.Character.Class)
			}

		case ChoicesMsg:
			m.Choices = []FeatureWithChoices(msg)
			if len(m.Choices) > m.ChoiceIndex {
				m.ActiveFeature = m.Choices[m.ChoiceIndex].Name
				m.ActiveChoices = m.Choices[m.ChoiceIndex].Choices
				for index, choice := range m.ActiveChoices {
					m.List.InsertItem(index, item(choice))
				}
			}

			// If there's no choices to make then submit the character, otherwise proceed to feature selection.
			if len(m.Choices) <= m.ChoiceIndex {
				return m, tea.Sequence(submitCharacter(m.Character), switchState(promptSelection))
			} else {
				return m, tea.Sequence(setupInitialCreateList("Choose the feature you'd like!", &m.ActiveChoices), switchState(promptSelection))
			}
		}
	case promptSelection:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					// Append the selected class feature to an array and move to the next index of the choices array
					m.Selections = append(m.Selections, FeatureWithChoices{m.ActiveFeature, []string{string(i)}})
					m.ChoiceIndex++

					// If we've reached the end of the choices array then submit the selections
					if len(m.Choices) <= m.ChoiceIndex {
						return m, tea.Sequence(
							tea.Batch(
								submitChoices(m.Character.Name, m.Selections),
								SubmitCharacter(m.Character),
							),
							switchState(characterCreated),
						)
					}

					// Otherwise update the feature and choices being decided on
					m.ActiveFeature = m.Choices[m.ChoiceIndex].Name
					m.ActiveChoices = m.Choices[m.ChoiceIndex].Choices

					// Empty the list component's items
					for len(m.List.Items()) > 0 {
						m.List.RemoveItem(0)
					}
					for index, choice := range m.ActiveChoices {
						m.List.InsertItem(index, item(choice))
					}
				}

				return m, nil
			}
		}
	case characterCreated:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "tab":
				m.ChoiceIndex = 0
				return m, tea.Sequence(switchState(promptName), switchParentState(showHome))
			case "esc":
				return m, tea.Quit
			}
		}
	}

	if slices.Contains([]State{promptName, promptClass, promptStats, promptItems, promptItems}, m.State) {
		m.TextInput, cmd = m.TextInput.Update(msg)
	} else {
		m.List, cmd = m.List.Update(msg)
	}

	return m, cmd
}
