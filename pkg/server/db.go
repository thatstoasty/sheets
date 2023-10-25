package server

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	// "github.com/BurntSushi/toml"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func readCSVData(fileName string) [][]string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// ReadAll reads all the records from the CSV file and Returns them as slice of slices of string and an error if any
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading records")
	}

	return records
}

func SetupDB() {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Class{})
	db.AutoMigrate(&ClassFeature{})
	db.AutoMigrate(&Race{})
	db.AutoMigrate(&Feat{})
	db.AutoMigrate(&Item{})
	db.AutoMigrate(&Weapon{})
	db.AutoMigrate(&Gear{})
	db.AutoMigrate(&Spell{})
	db.AutoMigrate(&Option{})
	db.AutoMigrate(&Character{})

	// Create classes
	classes := readCSVData("data/class.csv")
	for _, class := range classes {
		db.Save(&Class{Name: class[0]})
	}

	// Create Races
	races := readCSVData("data/race.csv")
	for _, race := range races {
		db.Save(&Race{Name: race[0], Options: race[1]})
	}

	// Create Feats
	feats := readCSVData("data/feat.csv")
	for _, feat := range feats {
		db.Save(&Feat{Name: feat[0], Options: feat[1]})
	}

	// Create Items
	items := readCSVData("data/item.csv")
	for _, item := range items {
		quantity, err := strconv.ParseInt(item[1], 10, 32)
		if err != nil {
			panic(err)
		}
		db.Save(&Item{Name: item[0], Quantity: uint16(quantity), Description: item[2], Properties: item[3], Options: item[4]})
	}

	// Create Weapons
	weapons := readCSVData("data/weapon.csv")
	for _, weapon := range weapons {
		db.Save(&Weapon{Name: weapon[0], Type: weapon[1], Description: weapon[2], Properties: weapon[3], Options: weapon[4]})
	}

	// Create Armor
	armors := readCSVData("data/armor.csv")
	for _, armor := range armors {
		db.Save(&Gear{Name: armor[0], Type: armor[1], Description: armor[2], Properties: armor[3], Options: armor[4]})
	}

	// Create Options
	options := readCSVData("data/option.csv")
	for _, option := range options {
		db.Save(&Option{Name: option[0], Type: option[1], Description: option[2]})
	}

	// Create Spells
	spells := readCSVData("data/spell.csv")
	for _, spell := range spells {
		level, err := strconv.ParseInt(spell[1], 10, 32)
		if err != nil {
			panic(err)
		}
		db.Save(&Spell{Name: spell[0], Level: uint16(level), Description: spell[2], Options: spell[3]})
	}

	// Create Class Features
	classFeatures := readCSVData("data/class_features.csv")
	for _, classFeature := range classFeatures {
		level, err := strconv.ParseInt(classFeature[3], 10, 32)
		if err != nil {
			panic(err)
		}
		db.Save(&ClassFeature{Name: classFeature[0], Class: classFeature[1], SubClass: classFeature[2], Level: uint32(level), Options: classFeature[4]})
	}

	db.Save(&Character{
		Name:         "Example",
		Class:        "Fighter,1,",
		HP:           "1",
		Proficiency:  "1",
		Strength:     "1",
		Dexterity:    "1",
		Constitution: "1",
		Intelligence: "1",
		Wisdom:       "1",
		Charisma:     "1",
		Race:         "Dwarf",
		Feats:        "",
		Items:        "",
	})
}
