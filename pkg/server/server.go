package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ClassInfo struct {
	Class    string
	Level    string
	SubClass string
}

type CharacterInfo struct {
	Name         string
	Race         string
	HP           string
	Proficiency  string
	Strength     string
	Dexterity    string
	Constitution string
	Intelligence string
	Wisdom       string
	Charisma     string
	ClassInfo    []ClassInfo
	Actions      []string
	BonusActions []string
	Passives     []string
}

func GetIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func getCharacterOptions(db *gorm.DB, name string) CharacterInfo {
	var character Character
	db.First(&character, "name = ?", name)

	var raceRecords []Race
	_ = db.Find(&raceRecords, "name in ?", []string{character.Race, "Default"})

	items := strings.Split(character.Items, ",")
	var itemRecords []Item
	_ = db.Find(&itemRecords, "name in ?", items)

	feats := strings.Split(character.Feats, ",")
	var featRecords []Feat
	_ = db.Find(&featRecords, "name in ?", feats)

	var options []string

	for _, race := range raceRecords {
		raceOptions := strings.Split(race.Options, "|")
		options = append(options, raceOptions...)
	}

	for _, feat := range featRecords {
		featOptions := strings.Split(feat.Options, "|")
		options = append(options, featOptions...)
	}

	for _, item := range itemRecords {
		itemOptions := strings.Split(item.Options, "|")
		options = append(options, itemOptions...)
	}

	var optionRecords []Option
	_ = db.Find(&optionRecords, "name in ?", options)

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

	classes := strings.Split(character.Class, ";")
	var classInfo []ClassInfo
	for _, class := range classes {
		classConfig := strings.Split(class, ",")
		classInfo = append(classInfo, ClassInfo{Class: classConfig[0], Level: classConfig[1], SubClass: classConfig[2]})
	}

	characterInfo := CharacterInfo{
		Name:         character.Name,
		Race:         character.Race,
		HP:           character.HP,
		Proficiency:  character.Proficiency,
		Strength:     character.Strength,
		Dexterity:    character.Dexterity,
		Constitution: character.Constitution,
		Intelligence: character.Intelligence,
		Wisdom:       character.Wisdom,
		Charisma:     character.Charisma,
		ClassInfo:    classInfo,
		Actions:      actions,
		BonusActions: bonusActions,
		Passives:     passives,
	}

	return characterInfo
}

func GetCharacter(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	characterInfo := getCharacterOptions(db, c.QueryParam("name"))

	return c.Render(http.StatusOK, "character", characterInfo)
}

func GetCharacterNames(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	var names []string
	db.Table("characters").Select("name").Scan(&names)

	return c.Render(http.StatusOK, "drop_down", names)
}

func GetOptionDescription(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	name := c.Request().Header["Hx-Trigger-Name"][0]
	log.Println(name)

	var option Option
	db.Table("options").Select("description").Where("name = ?", name).First(&option)

	return c.Render(http.StatusOK, "description", option.Description)
}

func SubmitCharacter(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		panic(err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var hero Character
	fileDecoder := toml.NewDecoder(file)
	_, err = fileDecoder.Decode(&hero)
	if err != nil {
		panic(err)
	}

	db.Save(&Character{
		Name:         hero.Name,
		Class:        hero.Class,
		HP:           hero.HP,
		Proficiency:  hero.Proficiency,
		Strength:     hero.Strength,
		Dexterity:    hero.Dexterity,
		Constitution: hero.Constitution,
		Intelligence: hero.Intelligence,
		Wisdom:       hero.Wisdom,
		Charisma:     hero.Charisma,
		Race:         hero.Race,
		Feats:        hero.Feats,
		Items:        hero.Items,
	})

	characterInfo := getCharacterOptions(db, hero.Name)
	log.Println(characterInfo)

	return c.Render(http.StatusOK, "character", characterInfo)
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
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

	// This will initiate our template renderer
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	server.Renderer = t

	server.GET("/", GetIndex)
	server.File("/css/output.css", "css/output.css")

	//// character
	server.POST("/character", SubmitCharacter)
	server.GET("/character", GetCharacter)
	server.GET("/character/names", GetCharacterNames)
	server.GET("/option", GetOptionDescription)

	// Start server
	server.Logger.Fatal(server.Start(":1323"))
}
