package server

import (
	"encoding/csv"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
	db.AutoMigrate(&ClassFeature{})
	db.AutoMigrate(&Characteristic{})
	db.AutoMigrate(&Item{})
	db.AutoMigrate(&Spell{})
	db.AutoMigrate(&Option{})
	db.AutoMigrate(&Character{})

	// Create Races
	races := readCSVData("data/players_handbook/races.csv")
	for _, race := range races {
		db.Save(&Characteristic{Name: race[0], Options: race[1]})
	}

	// Create Feats
	feats := readCSVData("data/players_handbook/feats.csv")
	for _, feat := range feats {
		db.Save(&Characteristic{Name: feat[0], Options: feat[1]})
	}

	// Create Items
	items := readCSVData("data/players_handbook/items.csv")
	for _, item := range items {
		db.Save(&Item{Name: item[0], Type: item[1], Description: item[2], Properties: item[3], Options: item[4]})
	}

	// Create Weapons
	filepath.Walk("data/players_handbook/weapons/", func(path string, file os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("File Name: %s\n", file.Name())

		if !file.IsDir() {
			weapons := readCSVData(path)
			for _, weapon := range weapons {
				db.Save(&Item{Name: weapon[0], Type: weapon[1], Description: weapon[2], Properties: weapon[3], Options: weapon[4]})
			}
		}

		return nil
	})

	// Create Armor
	armors := readCSVData("data/players_handbook/armor.csv")
	for _, armor := range armors {
		db.Save(&Item{Name: armor[0], Type: armor[1], Description: armor[2], Properties: armor[3], Options: armor[4]})
	}

	// Create Options
	options := readCSVData("data/players_handbook/options/base.csv")
	for _, option := range options {
		db.Save(&Option{Name: option[0], Type: option[1], Description: option[2]})
	}
	spellOptions := readCSVData("data/players_handbook/options/spells.csv")
	for _, option := range spellOptions {
		db.Save(&Option{Name: option[0], Type: option[1], Description: option[2]})
	}
	featOptions := readCSVData("data/players_handbook/options/feats.csv")
	for _, option := range featOptions {
		db.Save(&Option{Name: option[0], Type: option[1], Description: option[2]})
	}

	filepath.Walk("data/players_handbook/options/classes", func(path string, file os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("File Name: %s\n", file.Name())

		if !file.IsDir() {
			options := readCSVData(path)
			for _, option := range options {
				db.Save(&Option{Name: option[0], Type: option[1], Description: option[2]})
			}
		}

		return nil
	})

	// Create Spells
	spells := readCSVData("data/players_handbook/spells.csv")
	for _, spell := range spells {
		level, err := strconv.ParseInt(spell[1], 10, 32)
		if err != nil {
			panic(err)
		}
		db.Save(&Spell{Name: spell[0], Level: uint16(level), Description: spell[2], Options: spell[3]})
	}

	// Create Class Features
	filepath.Walk("data/players_handbook/classes", func(path string, file os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("File Name: %s\n", file.Name())

		if !file.IsDir() {
			classFeatures := readCSVData(path)
			for _, classFeature := range classFeatures {
				level, err := strconv.ParseInt(classFeature[3], 10, 32)
				if err != nil {
					panic(err)
				}
				db.Save(&ClassFeature{Name: classFeature[0], Class: classFeature[1], SubClass: classFeature[2], Level: uint32(level), Options: classFeature[4]})
			}
		}

		return nil
	})

	db.Save(&Character{
		Name:         "Example",
		Class:        "Fighter,1,",
		HP:           1,
		Proficiency:  1,
		Strength:     1,
		Dexterity:    1,
		Constitution: 1,
		Intelligence: 1,
		Wisdom:       1,
		Charisma:     1,
		Race:         "Dwarf",
		Feats:        "",
		Items:        "",
	})
}
