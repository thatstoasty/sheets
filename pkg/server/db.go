package server

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func SetupDB() {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&ClassFeature{})
	db.AutoMigrate(&Characteristic{})
	db.AutoMigrate(&FeatureChoices{})
	db.AutoMigrate(&Item{})
	db.AutoMigrate(&Weapon{})
	db.AutoMigrate(&Spell{})
	db.AutoMigrate(&Option{})
	db.AutoMigrate(&Character{})
	db.AutoMigrate(&Race{})

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
