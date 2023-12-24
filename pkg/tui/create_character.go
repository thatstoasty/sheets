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

	"github.com/thatstoasty/character-sheet-ui/pkg/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setStats(input string, character *Character) {
	stats := strings.Split(input, ",")
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

	character.HP = uint8(hp)
	character.Strength = uint8(strength)
	character.Dexterity = uint8(dexterity)
	character.Constitution = uint8(constitution)
	character.Intelligence = uint8(intelligence)
	character.Wisdom = uint8(wisdom)
	character.Charisma = uint8(charisma)
}

func submitCharacter(hero Character) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		db.Save(&Character{
			Name:           hero.Name,
			Class:          hero.Class,
			HP:             hero.HP,
			Strength:       hero.Strength,
			Dexterity:      hero.Dexterity,
			Constitution:   hero.Constitution,
			Intelligence:   hero.Intelligence,
			Wisdom:         hero.Wisdom,
			Charisma:       hero.Charisma,
			Race:           hero.Race,
			Feats:          hero.Feats,
			Items:          hero.Items,
			Helmet:         hero.Helmet,
			Cloak:          hero.Cloak,
			Jewelery1:      hero.Jewelery1,
			Jewelery2:      hero.Jewelery2,
			Jewelery3:      hero.Jewelery3,
			Boots:          hero.Boots,
			Gloves:         hero.Gloves,
			MainHandWeapon: hero.MainHandWeapon,
			OffHandWeapon:  hero.OffHandWeapon,
		})

		return SwitchStateMsg(characterCreated)
	}
}

type FeatureWithChoices struct {
	Name    string
	Choices []string
}

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

func getTableData(category string) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		var choices []string
		switch category {
		case "Race":
			db.Table("races").Select("name").Scan(&choices)
		case "Weapon":
			db.Table("items").Select("name").Where("category = 'Weapon'").Scan(&choices)
		case "LightWeapon":
			db.Table("items").Select("name").Where("category = 'Weapon' and properties like '%Light%'").Scan(&choices)
		default:
			db.Table("items").Select("name").Where("type = ?", category).Scan(&choices)
		}

		return TableDataMsg(TableData{category, choices})
	}
}

func submitChoices(character string, choices []FeatureWithChoices) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		for _, feature := range choices {
			_ = db.Save(&server.FeatureChoices{Character: character, Feature: feature.Name, Choice: feature.Choices[0]})
		}

		return nil
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

func (m CreateCharacterModel) Init() tea.Cmd {
	return tea.Batch(
		getTableData("Race"),
		setupTextInput,
	)
}

func (m CreateCharacterModel) View() string {
	if slices.Contains([]State{promptName, promptClass, promptStats, promptFeats, promptItems, promptItems, characterCreated}, m.State) {
		var prompt string
		switch m.State {
		case promptName:
			prompt = "Enter the name of the character:"
		case promptClass:
			prompt = "Enter the classes of the character (| and , delimited like this Paladin|10|Oath of Ancients,Monk|10|Way of the Shadow):"
		case promptStats:
			prompt = "Enter the stats of the character separated by commas (HP, Strength, Dexterity, Constitution, Intelligence, Wisdom, Charisma):"
		case promptFeats:
			prompt = "Enter the feats your character has separated by commas:"
		case promptItems:
			prompt = "Enter the items your character has separated by commas:"
		case characterCreated:
			return fmt.Sprintf(
				"%s has been created!\n\n%s\n\n",
				m.Character.Name,
				"(esc to quit, tab to return home)",
			) + "\n"
		}

		return centeredStyle.Render(
			fmt.Sprintf("%s\n\n%s\n\n%s",
				prompt,
				m.TextInput.View(),
				"(esc to quit)",
			),
		)
	} else {
		return centeredStyle.Render(m.List.View())
	}
}

func (m CreateCharacterModel) Update(msg tea.Msg) (CreateCharacterModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case ListMsg:
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
				return m, tea.Sequence(setupList("Choose your race", &m.Races), switchState(promptRace))
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
					getTableData("Cloak"),
					getTableData("Helmet"),
					getTableData("Armor"),
					getTableData("Jewelery"),
					getTableData("Boots"),
					getTableData("Gloves"),
					getTableData("Weapon"),
					getTableData("LightWeapon"),
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
				setStats(m.TextInput.Value(), &m.Character)
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
				return m, tea.Sequence(setupList("What helmet does your character have equipped?", &m.Helmets), switchState(promptHelmet))
			}
		}
	case promptHelmet:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Helmet = m.TextInput.Value()
				return m, tea.Sequence(setupList("What cloak does your character have equipped?", &m.Cloaks), switchState(promptCloak))
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
				return m, tea.Sequence(setupList("What armor does your character have equipped?", &m.Armor), switchState(promptArmor))
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
				return m, tea.Sequence(setupList("What's the first piece jewelery your character has equipped?", &m.Jewelery), switchState(promptJewelery1))
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
				return m, tea.Sequence(setupList("What's the second piece jewelery your character has equipped?", &m.Jewelery), switchState(promptJewelery2))
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
				return m, tea.Sequence(setupList("What's the third piece jewelery your character has equipped?", &m.Jewelery), switchState(promptJewelery3))
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
				return m, tea.Sequence(setupList("What boots does your character have equipped?", &m.Boots), switchState(promptBoots))
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
				return m, tea.Sequence(setupList("What gloves does your character have equipped?", &m.Gloves), switchState(promptGloves))
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
				return m, tea.Sequence(setupList("What main hand weapon does your character have equipped?", &m.Weapons), switchState(promptMainHandWeapon))
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
				return m, tea.Sequence(setupList("What off hand weapon does your character have equipped?", &m.LightWeapons), switchState(promptOffhandWeapon))
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
				return m, tea.Sequence(submitCharacter(m.Character))
			} else {
				return m, tea.Sequence(setupList("Choose the feature you'd like!", &m.ActiveChoices), switchState(promptSelection))
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
								submitCharacter(m.Character),
							),
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
