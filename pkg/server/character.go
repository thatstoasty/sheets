package server

import (
	"log"
	"math"
	"strconv"
	"strings"

	"gorm.io/gorm"
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
	var raceRecords []Characteristic
	_ = db.Find(&raceRecords, "name in ?", []string{characterRace, "Default"})

	for _, race := range raceRecords {
		raceOptions := strings.Split(race.Options, "|")
		*options = append(*options, raceOptions...)
	}
}

func addGearOptions(db *gorm.DB, character Character, options *[]string) {
	gear := []string{character.Helmet, character.Armor, character.Boots, character.Cloak, character.Jewelery1, character.Jewelery2, character.Jewelery3}
	var gearRecords []Characteristic
	_ = db.Find(&gearRecords, "name in ?", gear)

	for _, gear := range gearRecords {
		gearOptions := strings.Split(gear.Options, "|")
		*options = append(*options, gearOptions...)
	}

	weapons := []string{character.MainHandWeapon, character.OffHandWeapon}
	var weaponRecords []Characteristic
	_ = db.Find(&weaponRecords, "name in ?", weapons)

	for _, weapon := range gearRecords {
		weaponOptions := strings.Split(weapon.Options, "|")
		*options = append(*options, weaponOptions...)
	}
}

func addItemOptions(db *gorm.DB, characterItems string, options *[]string) {
	items := strings.Split(characterItems, ",")
	var itemRecords []Item
	_ = db.Find(&itemRecords, "name in ?", items)

	for _, item := range itemRecords {
		itemOptions := strings.Split(item.Options, "|")
		*options = append(*options, itemOptions...)
	}
}

func addFeatOptions(db *gorm.DB, characterFeats string, options *[]string) {
	feats := strings.Split(characterFeats, ",")
	var featRecords []Item
	_ = db.Find(&featRecords, "name in ?", feats)

	for _, feat := range featRecords {
		featOptions := strings.Split(feat.Options, "|")
		*options = append(*options, featOptions...)
	}
}

func addClassOptions(db *gorm.DB, characterClass string, characterLevel uint8, classInfo []ClassInfo, options *[]string) {
	allClassInfo := strings.Split(characterClass, ",")
	var allClassFeatureRecords []ClassFeature

	for _, class := range allClassInfo {
		classConfig := strings.Split(class, "|")

		var classFeatureRecords []ClassFeature
		_ = db.Where("class = ? AND level <= ? AND sub_class in ?", classConfig[0], classConfig[1], []string{"Base", classConfig[2]}).Find(&classFeatureRecords)
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

func GetCharacterOptions(db *gorm.DB, character Character, characterLevel uint8, classInfo []ClassInfo) []string {
	var options []string

	addRaceOptions(db, character.Race, &options)
	addGearOptions(db, character, &options)
	addItemOptions(db, character.Items, &options)
	addFeatOptions(db, character.Feats, &options)
	addClassOptions(db, character.Class, characterLevel, classInfo, &options)

	return options
}

func GetCharacterInfo(db *gorm.DB, name string) CharacterInfo {
	var character Character
	db.First(&character, "name = ?", name)

	var classInfo []ClassInfo
	var characterLevel uint8 = 0
	options := GetCharacterOptions(db, character, characterLevel, classInfo)

	var optionRecords []Option
	_ = db.Find(&optionRecords, "name in ?", options)

	var actions []Option
	var bonusActions []Option
	var passives []Option
	var reactions []Option
	var freeActions []Option
	var nonCombatActions []Option

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
