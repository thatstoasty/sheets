package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSheetsDirectory() {
	// Get user's home directory and create a .sheets directory if it doesn't exist
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/.sheets/file.db", home)); os.IsNotExist(err) {
		os.MkdirAll(fmt.Sprintf("%s/.sheets", home), 0700)
	}

	os.Setenv("SHEETS_DATABASE", fmt.Sprintf("%s/.sheets/file.db", home))
}

func getDatabaseSession() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(os.Getenv("SHEETS_DATABASE")), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	return db
}

func seedDatabase() {
	setupSheetsDirectory()
	db := getDatabaseSession()

	// Migrate the schema
	db.AutoMigrate(
		&ClassFeature{},
		&Characteristic{},
		&FeatureChoices{},
		&Item{},
		&Spell{},
		&Option{},
		&Character{},
	)

	// Create characterstics
	for _, characteristics := range characteristicsData {
		for _, characteristic := range characteristics {
			db.Save(&characteristic)
		}
	}

	// Create items
	for _, items := range itemsData {
		for _, item := range items {
			db.Save(&item)
		}
	}

	// Create Options
	for _, options := range optionsData {
		for _, option := range options {
			db.Save(&option)
		}
	}

	// Create Class Features
	for _, features := range classFeatures {
		for _, feature := range features {
			db.Save(&feature)
		}
	}

	// Create Spells
	for _, spell := range spellData {
		db.Save(&spell)
	}

	db.Save(&Character{
		Name:         "Example",
		Class:        "Fighter|1|",
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

func getCharacter(db *gorm.DB, name string) Character {
	var character Character
	db.Table("characters").Where("name = ?", name).First(&character)

	return character
}

func deleteCharacter(db *gorm.DB, name string) {
	db.Table("characters").Where("name = ?", name).Delete(&Character{})
}

func getCharacterNames(db *gorm.DB) []string {
	var names []string
	db.Table("characters").Select("name").Scan(&names)

	return names
}

func getRaceRecords(db *gorm.DB, race string) []Characteristic {
	var records []Characteristic
	db.Table("characteristics").Where("name in ? and type = 'Race'", []string{race, "Default"}).Scan(&records)

	return records
}

func getItemRecords(db *gorm.DB, items []string) []Item {
	var records []Item
	db.Table("items").Where("name in ?", items).Scan(&records)

	return records
}

func getItemCategoryRecords(db *gorm.DB, category string, items []string) []Item {
	var records []Item
	db.Table("items").Where("name in ? and category = ?", items, category).Scan(&records)

	return records
}

func getFeatRecords(db *gorm.DB, feats []string) []Characteristic {
	var records []Characteristic
	db.Table("characteristics").Where("name in ? and type = 'Feat'", feats).Scan(&records)

	return records
}

func getClassFeatureRecords(db *gorm.DB, class string, level string, subclass string) []ClassFeature {
	var records []ClassFeature
	db.Table("class_features").Where("class = ? and level = ? and sub_class in ?", class, level, []string{"Base", subclass}).Scan(&records)

	return records
}

func getOptionRecords(db *gorm.DB, options []string) []Option {
	var records []Option
	db.Table("options").Where("name in ?", options).Scan(&records)

	return records
}
