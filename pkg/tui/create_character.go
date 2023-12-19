package tui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

type RacesMsg []string

func getRaces() tea.Msg {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	var races []string
	db.Table("races").Select("name").Scan(&races)
	return RacesMsg(races)
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

type SetupTextInputMsg struct {
	TextInput textinput.Model
	ModelName string
}

func setupTextInput(modelName string) tea.Cmd {
	return func() tea.Msg {
		return SetupTextInputMsg{BuildTextInput(), modelName}
	}
}

func (m CreateCharacterModel) Init() tea.Cmd {
	return tea.Batch(
		getRaces,
		getCharacterNames,
		setupTextInput("CreateCharacter"),
	)

}

func (m CreateCharacterModel) View() string {
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
	case promptHelmet:
		return fmt.Sprintf(
			"What helmet does your character have equipped? Leave blank if none:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case promptCloak:
		return fmt.Sprintf(
			"What cloak does your character have equipped? Leave blank if none:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case promptArmor:
		return fmt.Sprintf(
			"What armor does your character have equipped? Leave blank if none:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case promptJewelery1:
		return fmt.Sprintf(
			"What is the first piece of jewelery your character has equipped? Leave blank if none:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case promptJewelery2:
		return fmt.Sprintf(
			"What is the second piece of jewelery your character has equipped? Leave blank if none:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case promptJewelery3:
		return fmt.Sprintf(
			"What is the third piece of jewelery your character has equipped? Leave blank if none:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case promptBoots:
		return fmt.Sprintf(
			"What boots does your character have equipped? Leave blank if none:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case promptGloves:
		return fmt.Sprintf(
			"What gloves does your character have equipped? Leave blank if none:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case promptMainHandWeapon:
		return fmt.Sprintf(
			"What main hand or two handed weapon does your character have equipped? Leave blank if none:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case promptOffhandWeapon:
		return fmt.Sprintf(
			"What off hand weapon does your character have equipped? Leave blank if none:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
	case promptSelection:
		return "\n" + m.List.View()
	default:
		return "Unrecoverable state!"
	}

}

func (m CreateCharacterModel) Update(msg tea.Msg) (CreateCharacterModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {

		case "/":
			fmt.Println(m.Character)
			fmt.Println(m.Choices)
			fmt.Println(m.State)
			fmt.Println(m.Races)
			fmt.Println(m.TextInput)

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

	case RacesMsg:
		m.Races = []string(msg)

	case CreateListMsg:
		m.List = list.Model(msg)
	}

	switch m.State {
	case promptName:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Name = m.TextInput.Value()
				m.TextInput.Reset()
				m.State++
				return m, setupInitialCreateList("Choose your race", &m.Races)
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

				m.State++
			}
		}
	case promptClass:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Class = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
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
				m.State++
				m.TextInput.Reset()
			}
		}
	case promptFeats:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Feats = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
			}
		}
	case promptItems:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Items = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
			}
		}
	case promptHelmet:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Helmet = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
			}
		}
	case promptCloak:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Cloak = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
			}
		}
	case promptArmor:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Armor = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
			}
		}
	case promptJewelery1:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Jewelery1 = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
			}
		}
	case promptJewelery2:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Jewelery2 = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
			}
		}
	case promptJewelery3:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Jewelery3 = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
			}
		}
	case promptBoots:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Boots = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
			}
		}
	case promptGloves:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Gloves = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
			}
		}
	case promptMainHandWeapon:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.MainHandWeapon = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
			}
		}
	case promptOffhandWeapon:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.OffHandWeapon = m.TextInput.Value()
				m.State++
				m.TextInput.Reset()
			}
		}
		return m, getChoices(m.Character.Class)
	case promptSelection:
		switch msg := msg.(type) {
		case SubmitCharacterMsg:
			return m, SubmitCharacter(m.Character)

		case ChoicesMsg:
			m.Choices = []FeatureWithChoices(msg)
			if len(m.Choices) > m.ChoiceIndex {
				m.ActiveFeature = m.Choices[m.ChoiceIndex].Name
				m.ActiveChoices = m.Choices[m.ChoiceIndex].Choices
				for index, choice := range m.ActiveChoices {
					m.List.InsertItem(index, item(choice))
				}
			}

			return m, nil

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
						return m, submitChoices(m.Character.Name, m.Selections)
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
	}

	// if m.State == promptSelection {
	// 	m.List, cmd = m.List.Update(msg)
	// } else {
	// 	// TextInput component Update -> CreateCharacter component Update -> Parent Update. Bubbles up
	// 	m.TextInput, cmd = m.TextInput.Update(msg)
	// }
	if m.State == promptRace {
		m.List, cmd = m.List.Update(msg)
	} else {
		m.TextInput, cmd = m.TextInput.Update(msg)
	}
	// cmds = append(cmds, cmd)
	// return m, tea.Batch(cmds...)
	return m, cmd
}
