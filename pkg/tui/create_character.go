package tui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
			db.Table("characteristics").Distinct("name").Where("type = ?", "Race").Scan(&choices)
		case "Weapon":
			db.Table("items").Select("name").Where("category = 'Weapon'").Scan(&choices)
		case "LightWeapon":
			db.Table("items").Select("name").Where("category = 'Weapon' and properties like '%Light%'").Scan(&choices)
		case "Class":
			db.Table("class_features").Distinct("class").Scan(&choices)
		case "Feat":
			db.Table("characteristics").Distinct("name").Where("type = ?", "Feat").Scan(&choices)
		default:
			db.Table("items").Select("name").Where("type = ?", category).Scan(&choices)
		}

		return TableDataMsg(TableData{category, choices})
	}
}

func getSubclasses(class string) tea.Cmd {
	return func() tea.Msg {
		db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect database")
		}

		var choices []string
		db.Table("class_features").Distinct("sub_class").Where("class = ?", class).Scan(&choices)

		return TableDataMsg(TableData{"Subclass", choices})
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
	promptTotalLevel
	promptClass
	promptLevel
	promptSubclass
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

type DNDData struct {
	Races        []string
	Weapons      []string
	LightWeapons []string
	Armor        []string
	Cloaks       []string
	Helmets      []string
	Boots        []string
	Gloves       []string
	Jewelery     []string
	Classes      []string
	Subclasses   []string
	Feats        []string
}

type Class struct {
	Name     string
	Level    int
	Subclass string
}

type CreateCharacterModel struct {
	State         State
	TextInput     textinput.Model
	Character     Character
	Classes       []Class
	Class         string
	Level         int
	Subclass      string
	TotalLevel    int
	List          list.Model
	Choices       []FeatureWithChoices
	ChoiceIndex   int
	ActiveFeature string
	ActiveChoices []string
	Selections    []FeatureWithChoices
	DND           DNDData
}

func (m CreateCharacterModel) Init() tea.Cmd {
	return tea.Batch(
		getTableData("Race"),
		setupTextInput,
	)
}

func (m CreateCharacterModel) characterView() string {
	if m.Character.Name == "" {
		return characterStyle.Render("You'll see more here after giving your character a name!\n\n")
	} else {
		// Format classes into a string for display
		var classes []string
		for _, class := range m.Classes {
			classes = append(classes, fmt.Sprintf("Level %d %s - %s", class.Level, class.Name, class.Subclass))
		}

		// Render the string that represents the character view
		return characterStyle.Render(lipgloss.JoinVertical(lipgloss.Top,
			fmt.Sprintf("%s: The level %d %s", m.Character.Name, m.TotalLevel, m.Character.Race),
			"------------",
			fmt.Sprintf("%d   %d   %d   %d   %d   %d   %d", m.Character.HP, m.Character.Strength, m.Character.Dexterity, m.Character.Constitution, m.Character.Intelligence, m.Character.Wisdom, m.Character.Charisma),
			"HP  STR  DEX  CON  INT  WIS  CHA\n",
			fmt.Sprintf("Class    | %s", strings.Join(classes, ", ")),
			fmt.Sprintf("Feats    | %s", m.Character.Feats),
			fmt.Sprintf("Items    | %s", m.Character.Items),
			fmt.Sprintf("Armor    | %s, %s, %s, %s, %s", m.Character.Helmet, m.Character.Cloak, m.Character.Armor, m.Character.Boots, m.Character.Gloves),
			fmt.Sprintf("Jewelery | %s, %s, %s", m.Character.Jewelery1, m.Character.Jewelery2, m.Character.Jewelery3),
			fmt.Sprintf("Weapons  | %s, %s", m.Character.MainHandWeapon, m.Character.OffHandWeapon),
		))
	}
}

func (m CreateCharacterModel) promptView() string {
	if m.State == characterCreated {
		return fmt.Sprintf(
			"%s has been created!\n\n%s\n\n",
			"Your character",
			"(esc to quit, tab to return home)",
		) + "\n"
	}

	if slices.Contains([]State{promptName, promptTotalLevel, promptLevel, promptStats, promptItems}, m.State) {
		var prompt string
		switch m.State {
		case promptName:
			prompt = "Enter the name of the character:"
		case promptTotalLevel:
			prompt = "Enter your character's total level:"
		case promptLevel:
			prompt = "Enter the level of the class:"
		case promptStats:
			prompt = "Enter the stats of the character separated by commas\n(HP, Strength, Dexterity, Constitution, Intelligence, Wisdom, Charisma):"
		case promptItems:
			prompt = "Enter the items your character has separated by commas:"
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

func (m CreateCharacterModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.promptView(), m.characterView())
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
			m.DND.Weapons = tableData.Records
		case "LightWeapon":
			m.DND.LightWeapons = tableData.Records
		case "Armor":
			m.DND.Armor = tableData.Records
		case "Helmet":
			m.DND.Helmets = tableData.Records
		case "Cloak":
			m.DND.Cloaks = tableData.Records
		case "Boots":
			m.DND.Boots = tableData.Records
		case "Gloves":
			m.DND.Gloves = tableData.Records
		case "Jewelery":
			m.DND.Jewelery = tableData.Records
		case "Race":
			m.DND.Races = tableData.Records
		case "Class":
			m.DND.Classes = tableData.Records
		case "Subclass":
			m.DND.Subclasses = tableData.Records
		case "Feat":
			m.DND.Feats = tableData.Records
		}
	}

	switch m.State {
	case promptName:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Name = m.TextInput.Value()
				return m, tea.Sequence(setupList("Choose your race", &m.DND.Races), switchState(promptRace))
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
					getTableData("Class"),
					getTableData("Feat"),
					switchState(promptTotalLevel),
				)
			}
		}
	case promptTotalLevel:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				level, err := strconv.Atoi(m.TextInput.Value())
				if err != nil {
					log.Fatal("Failed to convert level to an integer.")
				}

				if level < 1 || level > 20 {
					log.Fatal("Level must be between 1 and 20.")
				}
				m.TotalLevel = level
				return m, tea.Sequence(setupList("Choose your class:", &m.DND.Classes), switchState(promptClass))
			}
		}
	case promptClass:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Class = string(i)
				}

				return m, tea.Batch(getSubclasses(m.Class), switchState(promptLevel))
			}
		}
	case promptLevel:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				level, err := strconv.Atoi(m.TextInput.Value())
				if err != nil {
					log.Fatal("Failed to convert level to an integer.")
				}
				m.Level = level

				return m, tea.Sequence(setupList("Choose your subclass:", &m.DND.Subclasses), switchState(promptSubclass))
			}
		}
	case promptSubclass:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Subclass = string(i)
				}

				// Add submitted class to list of classes
				m.Classes = append(m.Classes, Class{m.Class, m.Level, m.Subclass})

				// Take a sum of the levels of the classes submitted so far
				var levelSum uint8
				for _, class := range m.Classes {
					levelSum += uint8(class.Level)
				}

				// Then check if the sum of the class levels is greater than or equal to the total level provided earlier.
				// If equal, move on to prompt stats otherwise throw an error.
				if levelSum == uint8(m.TotalLevel) {
					var classStrings []string
					for _, class := range m.Classes {
						classStrings = append(classStrings, fmt.Sprintf("%s|%d|%s", class.Name, class.Level, class.Subclass))
					}
					m.Character.Class = strings.Join(classStrings, ",")
					return m, switchState(promptStats)
				} else if levelSum < uint8(m.TotalLevel) {
					m.Class = ""
					m.Level = 0
					m.Subclass = ""
					return m, tea.Sequence(setupList("Choose your class:", &m.DND.Classes), switchState(promptClass))
				} else {
					log.Fatal("The sum of the class levels is greater than the total level provided!")
				}
			}
		}
	case promptStats:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				setStats(m.TextInput.Value(), &m.Character)
				return m, tea.Sequence(setupList("Choose your feats:", &m.DND.Feats), switchState(promptFeats))
			}
		}
	case promptFeats:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				i, ok := m.List.SelectedItem().(item)
				if ok {
					m.Character.Feats = string(i)
				}
				return m, switchState(promptItems)
			}
		}
	case promptItems:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Items = m.TextInput.Value()
				return m, tea.Sequence(setupList("What helmet does your character have equipped?", &m.DND.Helmets), switchState(promptHelmet))
			}
		}
	case promptHelmet:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "enter":
				m.Character.Helmet = m.TextInput.Value()
				return m, tea.Sequence(setupList("What cloak does your character have equipped?", &m.DND.Cloaks), switchState(promptCloak))
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
				return m, tea.Sequence(setupList("What armor does your character have equipped?", &m.DND.Armor), switchState(promptArmor))
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
				return m, tea.Sequence(setupList("What's the first piece jewelery your character has equipped?", &m.DND.Jewelery), switchState(promptJewelery1))
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
				return m, tea.Sequence(setupList("What's the second piece jewelery your character has equipped?", &m.DND.Jewelery), switchState(promptJewelery2))
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
				return m, tea.Sequence(setupList("What's the third piece jewelery your character has equipped?", &m.DND.Jewelery), switchState(promptJewelery3))
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
				return m, tea.Sequence(setupList("What boots does your character have equipped?", &m.DND.Boots), switchState(promptBoots))
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
				return m, tea.Sequence(setupList("What gloves does your character have equipped?", &m.DND.Gloves), switchState(promptGloves))
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
				return m, tea.Sequence(setupList("What main hand weapon does your character have equipped?", &m.DND.Weapons), switchState(promptMainHandWeapon))
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
				return m, tea.Sequence(setupList("What off hand weapon does your character have equipped?", &m.DND.LightWeapons), switchState(promptOffhandWeapon))
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
				m.Character = Character{}
				m.ChoiceIndex = 0
				return m, tea.Sequence(switchState(promptName), switchParentState(showHome))
			case "esc":
				return m, tea.Quit
			}
		}
	}

	if slices.Contains([]State{promptName, promptTotalLevel, promptLevel, promptStats, promptItems}, m.State) {
		m.TextInput, cmd = m.TextInput.Update(msg)
	} else {
		m.List, cmd = m.List.Update(msg)
	}

	return m, cmd
}
