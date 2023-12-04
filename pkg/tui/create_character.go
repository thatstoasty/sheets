package tui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	// "log"
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
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
)

type CreateCharacterModel struct {
	State     State
	TextInput textinput.Model
	Err       error
	Character Character
}

func BuildTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Mikhail"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return ti
}

// type setupMsg bool

// func initialTextInput() tea.Cmd {
// 	return func () tea.Msg {
// 		return setupMsg(true)
// 	}
// }

func (m CreateCharacterModel) Init() tea.Cmd {
	// return initialTextInput()
	return nil
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
		return fmt.Sprintf(
			"Enter the race of the character:\n\n%s\n\n%s",
			m.TextInput.View(),
			"(esc to quit)",
		) + "\n"
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
	default:
		return "Unrecoverable state!"
	}

}

func (m CreateCharacterModel) Update(msg tea.Msg) (CreateCharacterModel, tea.Cmd) {
	switch msg := msg.(type) {

	// case setupMsg:
	// 	switch msg {
	// 	case true:
	// 		m.TextInput = BuildTextInput()
	// 	}

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {

		case "/":
			fmt.Println(m.Character)

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":

			// Depending on the current prompt, save the value in the corresponding attribute of the character
			switch m.State {
			case promptName:
				m.Character.Name = m.TextInput.Value()
			case promptRace:
				m.Character.Race = m.TextInput.Value()
			case promptClass:
				m.Character.Class = m.TextInput.Value()
			case promptStats:
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

			case promptFeats:
				m.Character.Race = m.TextInput.Value()
			case promptItems:
				m.Character.Items = m.TextInput.Value()
			case promptHelmet:
				m.Character.Helmet = m.TextInput.Value()
			case promptCloak:
				m.Character.Cloak = m.TextInput.Value()
			case promptArmor:
				m.Character.Armor = m.TextInput.Value()
			case promptJewelery1:
				m.Character.Jewelery1 = m.TextInput.Value()
			case promptJewelery2:
				m.Character.Jewelery2 = m.TextInput.Value()
			case promptJewelery3:
				m.Character.Jewelery3 = m.TextInput.Value()
			case promptBoots:
				m.Character.Boots = m.TextInput.Value()
			case promptGloves:
				m.Character.Gloves = m.TextInput.Value()
			case promptMainHandWeapon:
				m.Character.MainHandWeapon = m.TextInput.Value()
			case promptOffhandWeapon:
				m.Character.OffHandWeapon = m.TextInput.Value()
				return m, SubmitCharacter(m.Character)
			}
			m.State++
			m.TextInput = BuildTextInput()
			return m, nil
		}
	}

	// TextInput component Update -> CreateCharacter component Update -> Parent Update. Bubbles up
	var cmd tea.Cmd
	m.TextInput, cmd = m.TextInput.Update(msg)

	return m, cmd
}
