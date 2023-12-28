package server

import (
	"log"
	"math"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/thatstoasty/character-sheet-ui/pkg/database"
)

func CalculateAbilityScore(score uint8) int8 {
	abilityModifier := float64(score) - 10
	if abilityModifier == 0 {
		return int8(abilityModifier)
	} else {
		return int8(math.Floor(abilityModifier / 2))
	}
}

func CalculateProficiencyScore(characterLevel uint8) uint8 {
	switch {
	case characterLevel >= 17:
		return 6
	case characterLevel >= 13:
		return 5
	case characterLevel >= 9:
		return 4
	case characterLevel >= 5:
		return 3
	default:
		return 2
	}
}

func addRaceOptions(db *gorm.DB, characterRace string, options *[]string) {
	raceRecords := database.GetRaceRecords(db, characterRace)

	for _, race := range raceRecords {
		raceOptions := strings.Split(race.Options, "|")
		*options = append(*options, raceOptions...)
	}
}

func addGearOptions(db *gorm.DB, character database.Character, options *[]string) {
	armor := []string{character.Helmet, character.Armor, character.Boots, character.Cloak, character.Jewelery1, character.Jewelery2, character.Jewelery3}
	armorRecords := database.GetItemCategoryRecords(db, "Armor", armor)
	for _, gear := range armorRecords {
		gearOptions := strings.Split(gear.Options, "|")
		*options = append(*options, gearOptions...)
	}

	weapons := []string{character.MainHandWeapon, character.OffHandWeapon}
	weaponRecords := database.GetItemCategoryRecords(db, "Weapon", weapons)
	for _, weapon := range weaponRecords {
		weaponOptions := strings.Split(weapon.Options, "|")
		*options = append(*options, weaponOptions...)
	}
}

func addItemOptions(db *gorm.DB, characterItems string, options *[]string) {
	items := strings.Split(characterItems, ",")
	itemRecords := database.GetItemRecords(db, items)
	for _, item := range itemRecords {
		itemOptions := strings.Split(item.Options, "|")
		*options = append(*options, itemOptions...)
	}
}

func addFeatOptions(db *gorm.DB, characterFeats string, options *[]string) {
	feats := strings.Split(characterFeats, ",")
	featRecords := database.GetFeatRecords(db, feats)
	for _, feat := range featRecords {
		featOptions := strings.Split(feat.Options, "|")
		*options = append(*options, featOptions...)
	}
}

type ClassInfo struct {
	Class    string
	Level    string
	SubClass string
}

func addClassOptions(db *gorm.DB, characterClass string, characterLevel uint8, classInfo []ClassInfo, options *[]string) {
	allClassInfo := strings.Split(characterClass, ",")
	var allClassFeatureRecords []database.ClassFeature

	for _, class := range allClassInfo {
		classConfig := strings.Split(class, "|")

		classFeatureRecords := database.GetClassFeatureRecords(db, classConfig[0], classConfig[1], classConfig[2])
		allClassFeatureRecords = append(allClassFeatureRecords, classFeatureRecords...)
		classInfo = append(classInfo, ClassInfo{Class: classConfig[0], Level: classConfig[1], SubClass: classConfig[2]})

		classLevel, err := strconv.Atoi(classConfig[1])
		if err != nil {
			log.Fatal("Failed to convert class level to an integer.")
		}
		characterLevel += uint8(classLevel)
	}

	for _, classFeature := range allClassFeatureRecords {
		classFeatureOptions := strings.Split(classFeature.Options, "|")
		*options = append(*options, classFeatureOptions...)
	}
}

func GetCharacterOptions(db *gorm.DB, character database.Character, characterLevel uint8, classInfo []ClassInfo) []string {
	var options []string

	addRaceOptions(db, character.Race, &options)
	addGearOptions(db, character, &options)
	addItemOptions(db, character.Items, &options)
	addFeatOptions(db, character.Feats, &options)
	addClassOptions(db, character.Class, characterLevel, classInfo, &options)

	return options
}

type ScoreWithModifier struct {
	Score    uint8
	Modifier int8
}

type CharacterInfo struct {
	Name             string
	Race             string
	HP               uint8
	Proficiency      uint8
	Strength         ScoreWithModifier
	Dexterity        ScoreWithModifier
	Constitution     ScoreWithModifier
	Intelligence     ScoreWithModifier
	Wisdom           ScoreWithModifier
	Charisma         ScoreWithModifier
	ClassInfo        []ClassInfo
	Actions          []database.Option
	BonusActions     []database.Option
	Passives         []database.Option
	Reactions        []database.Option
	FreeActions      []database.Option
	NonCombatActions []database.Option
	Items            []string
	Equipped         []string
}

func GetCharacterInfo(db *gorm.DB, name string) CharacterInfo {
	character := database.GetCharacter(db, name)

	var (
		classInfo      []ClassInfo
		characterLevel uint8 = 0
	)
	options := GetCharacterOptions(db, character, characterLevel, classInfo)
	optionRecords := database.GetOptionRecords(db, options)

	var (
		actions          []database.Option
		bonusActions     []database.Option
		passives         []database.Option
		reactions        []database.Option
		freeActions      []database.Option
		nonCombatActions []database.Option
	)

	for _, opt := range optionRecords {
		switch {
		case opt.Type == "Action":
			actions = append(actions, opt)
		case opt.Type == "BonusAction":
			bonusActions = append(bonusActions, opt)
		case opt.Type == "Passive":
			passives = append(passives, opt)
		case opt.Type == "Reaction":
			reactions = append(reactions, opt)
		case opt.Type == "FreeAction":
			freeActions = append(freeActions, opt)
		case opt.Type == "NonCombatAction":
			nonCombatActions = append(nonCombatActions, opt)
		}
	}

	characterInfo := CharacterInfo{
		Name:             character.Name,
		Race:             character.Race,
		HP:               character.HP,
		Proficiency:      CalculateProficiencyScore(uint8(characterLevel)),
		Strength:         ScoreWithModifier{character.Strength, CalculateAbilityScore(character.Strength)},
		Dexterity:        ScoreWithModifier{character.Dexterity, CalculateAbilityScore(character.Dexterity)},
		Constitution:     ScoreWithModifier{character.Constitution, CalculateAbilityScore(character.Constitution)},
		Intelligence:     ScoreWithModifier{character.Intelligence, CalculateAbilityScore(character.Intelligence)},
		Wisdom:           ScoreWithModifier{character.Wisdom, CalculateAbilityScore(character.Wisdom)},
		Charisma:         ScoreWithModifier{character.Charisma, CalculateAbilityScore(character.Charisma)},
		ClassInfo:        classInfo,
		Actions:          actions,
		BonusActions:     bonusActions,
		Passives:         passives,
		Reactions:        reactions,
		FreeActions:      freeActions,
		NonCombatActions: nonCombatActions,
		Items:            strings.Split(character.Items, ","),
		Equipped:         []string{character.Helmet, character.Armor, character.Boots, character.Cloak, character.Jewelery1, character.Jewelery2, character.Jewelery3, character.MainHandWeapon, character.OffHandWeapon},
	}

	return characterInfo
}
