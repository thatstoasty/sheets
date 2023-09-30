package main

import (
	"encoding/json"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Metadata struct {
}

type Option struct {
	ID   uint32 `gorm:"primaryKey,autoIncrement"`
	Name string
	Type string
}

type Action struct {
	ID   uint32 `gorm:"primaryKey,autoIncrement"`
	Name string
}

type BonusAction struct {
	ID   uint32 `gorm:"primaryKey,autoIncrement"`
	Name string
}

type Passive struct {
	ID   uint32 `gorm:"primaryKey,autoIncrement"`
	Name string
}

type Class struct {
	ID   uint32 `gorm:"primaryKey,autoIncrement"`
	Name string
}

type ClassFeature struct {
	ID      uint32 `gorm:"primaryKey,autoIncrement"`
	Name    string
	Class   string
	Level   uint32
	Options string
}

type Race struct {
	ID      uint32 `gorm:"primaryKey,autoIncrement"`
	Name    string
	Options string
}

type Feat struct {
	ID      uint32 `gorm:"primaryKey,autoIncrement"`
	Name    string
	Options string
}

type Item struct {
	ID      uint32 `gorm:"primaryKey,autoIncrement"`
	Name    string
	Options string
}

type ClassLevel struct {
	Class string
	Level uint32
}

type Character struct {
	ID    uint32 `gorm:"primaryKey,autoIncrement"`
	Name  string
	Class string
	Race  string
	Feats string
	Items string
}

// type Spell struct {
// 	ID  		uint32 `gorm:"primaryKey,autoIncrement"`
// 	Name	 	string
// 	Description string
// 	Level		uint32
// }

func convertJSONString(jsonString string, targetObj any) {
	err := json.Unmarshal([]byte(jsonString), &targetObj)
	if err != nil {
		panic(err)
	}
}

func main() {
	// connect
	// db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
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

	// Create
	db.Create(&Class{Name: "Fighter"})

	// Read
	var class Class
	db.First(&class, 1)                     // find product with integer primary key
	db.First(&class, "name = ?", "Fighter") // find product with code D42
	log.Println(class)

	// Delete - delete product
	db.Delete(&class, 1)

	// Create classes
	db.Create(&Class{Name: "Artificer"})
	db.Create(&Class{Name: "Barbarian"})
	db.Create(&Class{Name: "Bard"})
	db.Create(&Class{Name: "Cleric"})
	db.Create(&Class{Name: "Druid"})
	db.Create(&Class{Name: "Fighter"})
	db.Create(&Class{Name: "Monk"})
	db.Create(&Class{Name: "Paladin"})
	db.Create(&Class{Name: "Ranger"})
	db.Create(&Class{Name: "Rogue"})
	db.Create(&Class{Name: "Sorcerer"})
	db.Create(&Class{Name: "Warlock"})
	db.Create(&Class{Name: "Wizard"})

	var classes []Class
	_ = db.Find(&classes)

	log.Println(classes)

	// Create Races
	db.Create(&Race{Name: "Dwarf", Options: `["Darkvision"]`})

	// Create Feats
	db.Create(&Feat{Name: "Warcaster", Options: `["Cast War"]`})
	db.Create(&Feat{Name: "Polearm Master", Options: `["Polearm Master Extra Attack"]`})

	// Create Items
	db.Create(&Feat{Name: "Rapier", Options: ""})

	// Create Actions
	db.Create(&Option{Name: "Cast War", Type: "Action"})

	// Create Bonus Actions
	db.Create(&Option{Name: "Polearm Master Extra Attack", Type: "Action"})

	// Create Passives
	db.Create(&Option{Name: "Darkvision", Type: "Action"})

	// var options []Option
	// _ = db.Find(&options)

	// log.Println(options)

	// var conf Config
	// _, err := toml.Decode(tomlData, &conf)

	// Create Character
	db.Create(&Character{
		Name:  "Mikhail",
		Class: `[{"Class": "Paladin", "Level": 1}]`,
		Race:  "Dwarf",
		Feats: `["Warcaster"]`,
		Items: `["Rapier"]`,
	})

	var character Character
	db.First(&character, "name = ?", "Mikhail")

	log.Println(character)
	var raceOptions Race
	_ = db.First(&raceOptions, "name = ?", character.Race)

	var items []string
	// convertJSONString(character.Items, items)
	err = json.Unmarshal([]byte(character.Items), &items)
	if err != nil {
		panic(err)
	}
	var itemOptions Item
	_ = db.First(&itemOptions, "name = ?", items[0])

	var feats []string
	err = json.Unmarshal([]byte(character.Feats), &feats)
	if err != nil {
		panic(err)
	}
	var featOptions Feat
	_ = db.First(&featOptions, "name = ?", feats[0])

	var featsOptions []Feat
	_ = db.Find(&featsOptions)
	log.Println(featsOptions)

	log.Println(raceOptions.Options)
	log.Println(itemOptions.Options)
	log.Println(featOptions.Options)

	var raceOpt []string
	err = json.Unmarshal([]byte(raceOptions.Options), &raceOpt)
	if err != nil {
		panic(err)
	}
	log.Println(raceOpt[0])

	// var itemOpt []string
	// err = json.Unmarshal([]byte(itemOptions.Options), &itemOpt)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(itemOpt)

	var featOpt []string
	err = json.Unmarshal([]byte(featOptions.Options), &featOpt)
	if err != nil {
		panic(err)
	}
	log.Println(featOpt[0])

	var options []string
	options = append(options, raceOpt...)
	// options = append(options, itemOpt...)
	options = append(options, featOpt...)

	log.Println(options)

	var optionRecords []Option
	_ = db.Find(&optionRecords, "name in ?", options)

	log.Println(optionRecords)
	var actions []string
	var bonusActions []string
	var passives []string

	for _, opt := range optionRecords {
		switch {
		case opt.Type == "Action":
			actions = append(actions, opt.Name)
		case opt.Type == "BonusAction":
			bonusActions = append(bonusActions, opt.Name)
		case opt.Type == "Passive":
			passives = append(passives, opt.Name)
		}
	}

	log.Println(actions)
	log.Println(bonusActions)
	log.Println(passives)

}
