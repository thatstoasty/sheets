package server

import (
	"log"

	// "github.com/BurntSushi/toml"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDB() {
	// connect
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	// db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
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

	// // Create
	// db.Create(&Class{Name: "Fighter"})

	// // Read
	// var class Class
	// db.First(&Class, 1)                     // find product with integer primary key
	// db.First(&Class, "name = ?", "Fighter") // find product with code D42
	// log.Println(class)

	// // Delete - delete product
	// db.Delete(&Class, 1)

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

	// var classes []Class
	// _ = db.Find(&Classes)

	// log.Println(classes)

	// Create Races
	db.Create(&Race{Name: "Dwarf", Options: `["Darkvision"]`})

	// Create Feats
	db.Create(&Feat{Name: "Warcaster", Options: `["Cast War"]`})
	db.Create(&Feat{Name: "Polearm Master", Options: `["Polearm Master Extra Attack"]`})

	// Create Items
	db.Create(&Item{Name: "Rapier", Options: ""})

	// Create Actions
	db.Create(&Option{Name: "Cast War", Type: "Action"})

	// Create Bonus Actions
	db.Create(&Option{Name: "Polearm Master Extra Attack", Type: "Action"})

	// Create Passives
	db.Create(&Option{Name: "Darkvision", Type: "Passive"})

	// // var options []Option
	// // _ = db.Find(&Options)

	// // log.Println(options)

	// // var hero Character
	// // _, err = toml.DecodeFile("character.toml", &hero)
	// // if err != nil {
	// //     log.Println(err)
	// //     return
	// // }
	// // log.Println(hero)

	// Create Character
	// db.Create(&Character{
	// 	Name:  "Mikhail",
	// 	Class: `[{"Class": "Paladin", "Level": 1}]`,
	// 	Race:  "Dwarf",
	// 	Feats: `["Warcaster"]`,
	// 	Items: `["Rapier"]`,
	// })

	// var character Character
	// db.First(&Character, "name = ?", "Mikhail")

	// log.Println(character)
	// var raceOptions Race
	// _ = db.First(&RaceOptions, "name = ?", character.Race)

	// var items []string
	// // convertJSONString(character.Items, items)
	// err = json.Unmarshal([]byte(character.Items), &items)
	// if err != nil {
	// 	panic(err)
	// }
	// var itemOptions Item
	// _ = db.First(&itemOptions, "name = ?", items[0])

	// var feats []string
	// err = json.Unmarshal([]byte(character.Feats), &Feats)
	// if err != nil {
	// 	panic(err)
	// }
	// var featOptions Feat
	// _ = db.First(&FeatOptions, "name = ?", feats[0])

	// var featsOptions []Feat
	// _ = db.Find(&FeatsOptions)
	// log.Println(featsOptions)

	// log.Println(raceOptions.Options)
	// log.Println(itemOptions.Options)
	// log.Println(featOptions.Options)

	// var raceOpt []string
	// err = json.Unmarshal([]byte(raceOptions.Options), &RaceOpt)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(raceOpt[0])

	// // var itemOpt []string
	// // err = json.Unmarshal([]byte(itemOptions.Options), &itemOpt)
	// // if err != nil {
	// // 	panic(err)
	// // }
	// // log.Println(itemOpt)

	// var featOpt []string
	// err = json.Unmarshal([]byte(featOptions.Options), &FeatOpt)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(featOpt[0])

	// var options []string
	// options = append(options, raceOpt...)
	// // options = append(options, itemOpt...)
	// options = append(options, featOpt...)

	// log.Println(options)

	// var optionRecords []Option
	// _ = db.Find(&OptionRecords, "name in ?", options)

	// log.Println(optionRecords)
	// var actions []string
	// var bonusActions []string
	// var passives []string

	// for _, opt := range optionRecords {
	// 	switch {
	// 	case opt.Type == "Action":
	// 		actions = append(actions, opt.Name)
	// 	case opt.Type == "BonusAction":
	// 		bonusActions = append(bonusActions, opt.Name)
	// 	case opt.Type == "Passive":
	// 		passives = append(passives, opt.Name)
	// 	}
	// }

	// log.Println(actions)
	// log.Println(bonusActions)
	// log.Println(passives)
}
