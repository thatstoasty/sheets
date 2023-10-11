package server

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

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
		db.Save(&Item{Name: item[0], Options: item[1]})
	}

	// Create Options
	options := readCSVData("data/option.csv")
	for _, option := range options {
		db.Save(&Option{Name: option[0], Type: option[1], Description: option[2]})
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
