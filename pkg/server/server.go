package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

type ResourceWithDescription struct {
	Name        string
	Description string
}

type ScoreWithModifier struct {
	Score    uint8
	Modifier int8
}

type CharacterInfo struct {
	Name             string
	Race             string
	HP               uint8
	Proficiency      uint8
	Strength         ScoreWithModifier
	Dexterity        ScoreWithModifier
	Constitution     ScoreWithModifier
	Intelligence     ScoreWithModifier
	Wisdom           ScoreWithModifier
	Charisma         ScoreWithModifier
	ClassInfo        []ClassInfo
	Actions          []Option
	BonusActions     []Option
	Passives         []Option
	Reactions        []Option
	FreeActions      []Option
	NonCombatActions []Option
	Items            []string
	Equipped         []string
}

func GetIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func GetCharacter(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	characterInfo := GetCharacterInfo(db, c.QueryParam("name"))

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

type OptionWithDescription struct {
	Name        string
	Description string
}

func GetOptionDescription(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	name := c.Request().Header["Hx-Trigger-Name"][0]

	var option Option
	db.Table("options").Select("description").Where("name = ?", name).First(&option)

	return c.Render(http.StatusOK, "description", OptionWithDescription{name, option.Description})
}

func GetItemDescription(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	name := c.Request().Header["Hx-Trigger-Name"][0]

	var option Option
	db.Table("items").Select("description").Where("name = ?", name).First(&option)

	return c.Render(http.StatusOK, "description", OptionWithDescription{name, option.Description})
}

func GetGearDescription(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	name := c.Request().Header["Hx-Trigger-Name"][0]

	var option Option
	db.Table("gears").Select("description").Where("name = ?", name).First(&option)

	return c.Render(http.StatusOK, "description", OptionWithDescription{name, option.Description})
}

func GetWeaponDescription(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	name := c.Request().Header["Hx-Trigger-Name"][0]

	var option Option
	db.Table("weapons").Select("description").Where("name = ?", name).First(&option)

	return c.Render(http.StatusOK, "description", OptionWithDescription{name, option.Description})
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
		Name:           hero.Name,
		Class:          hero.Class,
		HP:             hero.HP,
		Proficiency:    hero.Proficiency,
		Strength:       hero.Strength,
		Dexterity:      hero.Dexterity,
		Constitution:   hero.Constitution,
		Intelligence:   hero.Intelligence,
		Wisdom:         hero.Wisdom,
		Charisma:       hero.Charisma,
		Race:           hero.Race,
		Feats:          hero.Feats,
		Items:          hero.Items,
		Helmet:         hero.Helmet,
		Cloak:          hero.Cloak,
		Jewelery1:      hero.Jewelery1,
		Jewelery2:      hero.Jewelery2,
		Jewelery3:      hero.Jewelery3,
		Boots:          hero.Boots,
		Gloves:         hero.Gloves,
		MainHandWeapon: hero.MainHandWeapon,
		OffHandWeapon:  hero.OffHandWeapon,
	})

	characterInfo := GetCharacterInfo(db, hero.Name)

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
	server.File("/favicon.ico", "content/favicon.ico")
	server.File("/content/gun.png", "content/gun.png")

	//// character
	server.POST("/character", SubmitCharacter)
	server.GET("/character", GetCharacter)
	server.GET("/character/names", GetCharacterNames)
	server.GET("/option", GetOptionDescription)
	server.GET("/item", GetItemDescription)
	server.GET("/gear", GetGearDescription)
	server.GET("/weapon", GetWeaponDescription)

	// Start server
	server.Logger.Fatal(server.Start(":1323"))
}
