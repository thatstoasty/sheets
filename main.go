package main

import (
	"encoding/json"
	// "log"

	// "github.com/BurntSushi/toml"
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"

	"github.com/thatstoasty/character-sheet-ui/pkg/server"
)

func convertJSONString(jsonString string, targetObj any) {
	err := json.Unmarshal([]byte(jsonString), &targetObj)
	if err != nil {
		panic(err)
	}
}

func main() {
	// connect
	// db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	// // db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal("failed to connect database")
	// }

	// // Migrate the schema
	// db.AutoMigrate(&server.Class{})
	// db.AutoMigrate(&server.ClassFeature{})
	// db.AutoMigrate(&server.Race{})
	// db.AutoMigrate(&server.Feat{})
	// db.AutoMigrate(&server.Item{})
	// db.AutoMigrate(&server.Option{})
	// db.AutoMigrate(&server.Character{})

	// // Create
	// db.Create(&server.Class{Name: "Fighter"})

	// // Read
	// var class server.Class
	// db.First(&server.Class, 1)                     // find product with integer primary key
	// db.First(&server.Class, "name = ?", "Fighter") // find product with code D42
	// log.Println(class)

	// // Delete - delete product
	// db.Delete(&server.Class, 1)

	// // Create classes
	// db.Create(&server.Class{Name: "Artificer"})
	// db.Create(&server.Class{Name: "Barbarian"})
	// db.Create(&server.Class{Name: "Bard"})
	// db.Create(&server.Class{Name: "Cleric"})
	// db.Create(&server.Class{Name: "Druid"})
	// db.Create(&server.Class{Name: "Fighter"})
	// db.Create(&server.Class{Name: "Monk"})
	// db.Create(&server.Class{Name: "Paladin"})
	// db.Create(&server.Class{Name: "Ranger"})
	// db.Create(&server.Class{Name: "Rogue"})
	// db.Create(&server.Class{Name: "Sorcerer"})
	// db.Create(&server.Class{Name: "Warlock"})
	// db.Create(&server.Class{Name: "Wizard"})

	// var classes []server.Class
	// _ = db.Find(&server.Classes)

	// log.Println(classes)

	// // Create Races
	// db.Create(&server.Race{Name: "Dwarf", Options: `["Darkvision"]`})

	// // Create Feats
	// db.Create(&server.Feat{Name: "Warcaster", Options: `["Cast War"]`})
	// db.Create(&server.Feat{Name: "Polearm Master", Options: `["Polearm Master Extra Attack"]`})

	// // Create Items
	// db.Create(&server.Feat{Name: "Rapier", Options: ""})

	// // Create Actions
	// db.Create(&server.Option{Name: "Cast War", Type: "Action"})

	// // Create Bonus Actions
	// db.Create(&server.Option{Name: "Polearm Master Extra Attack", Type: "Action"})

	// // Create Passives
	// db.Create(&server.Option{Name: "Darkvision", Type: "Action"})

	// // var options []Option
	// // _ = db.Find(&server.Options)

	// // log.Println(options)

	// // var hero Character
	// // _, err = toml.DecodeFile("character.toml", &hero)
	// // if err != nil {
	// //     log.Println(err)
	// //     return
	// // }
	// // log.Println(hero)

	// // Create Character
	// db.Create(&server.Character{
	// 	Name:  "Mikhail",
	// 	Class: `[{"Class": "Paladin", "Level": 1}]`,
	// 	Race:  "Dwarf",
	// 	Feats: `["Warcaster"]`,
	// 	Items: `["Rapier"]`,
	// })

	// var character server.Character
	// db.First(&server.Character, "name = ?", "Mikhail")

	// log.Println(character)
	// var raceOptions server.Race
	// _ = db.First(&server.RaceOptions, "name = ?", character.Race)

	// var items []string
	// // convertJSONString(character.Items, items)
	// err = json.Unmarshal([]byte(character.Items), &items)
	// if err != nil {
	// 	panic(err)
	// }
	// var itemOptions server.Item
	// _ = db.First(&itemOptions, "name = ?", items[0])

	// var feats []string
	// err = json.Unmarshal([]byte(character.Feats), &server.Feats)
	// if err != nil {
	// 	panic(err)
	// }
	// var featOptions server.Feat
	// _ = db.First(&server.FeatOptions, "name = ?", feats[0])

	// var featsOptions []server.Feat
	// _ = db.Find(&server.FeatsOptions)
	// log.Println(featsOptions)

	// log.Println(raceOptions.Options)
	// log.Println(itemOptions.Options)
	// log.Println(featOptions.Options)

	// var raceOpt []string
	// err = json.Unmarshal([]byte(raceOptions.Options), &server.RaceOpt)
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
	// err = json.Unmarshal([]byte(featOptions.Options), &server.FeatOpt)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(featOpt[0])

	// var options []string
	// options = append(options, raceOpt...)
	// // options = append(options, itemOpt...)
	// options = append(options, featOpt...)

	// log.Println(options)

	// var optionRecords []server.Option
	// _ = db.Find(&server.OptionRecords, "name in ?", options)

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

	server.SetupDB()
	server.Start()
}
