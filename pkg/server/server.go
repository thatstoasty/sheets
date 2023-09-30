package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetCharacter(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	var character Character
	log.Println(c.Param("name"))
	db.First(&character, "name = ?", c.Param("name"))

	return c.JSON(http.StatusOK, character)
}

func SubmitCharacter(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

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

	return c.JSON(http.StatusOK, map[string]string{"message": "Character creation successful."})
}

func GetOptions(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	var character Character
	log.Println(c.Param("name"))
	db.First(&character, "name = ?", c.Param("name"))

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

	return c.JSON(http.StatusOK, actions)
}

func Start() {
	// Echo instance
	server := echo.New()

	// Middleware
	server.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.RequestID(),
	)

	server.HTTPErrorHandler = func(err error, c echo.Context) {
		// Take required information from error and context and send it to a service like New Relic
		fmt.Println(c.Path(), c.QueryParams(), err.Error())

		// Call the default handler to return the HTTP response
		server.DefaultHTTPErrorHandler(err, c)
	}

	//// Exercise
	server.GET("/character/:name", GetCharacter)
	server.POST("/character", SubmitCharacter)
	server.GET("/character/:name/options", GetOptions)

	// Start server
	server.Logger.Fatal(server.Start(":1323"))
}
