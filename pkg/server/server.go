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

type Options struct {
	Actions      []string `json:"actions"`
	BonusActions []string `json:"bonus_actions"`
	Passives     []string `json:"passives"`
}

func GetIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func getCharacterOptions(db *gorm.DB, name string) Options {
	var character Character
	db.First(&character, "name = ?", name)

	var raceRecords Race
	_ = db.First(&raceRecords, "name = ?", character.Race)

	items := strings.Split(character.Items, ",")
	var itemRecords []Item
	_ = db.Find(&itemRecords, "name in ?", items)

	feats := strings.Split(character.Feats, ",")
	var featRecords []Feat
	_ = db.Find(&featRecords, "name in ?", feats)

	var options []string
	raceOptions := strings.Split(raceRecords.Options, "|")
	options = append(options, raceOptions...)

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

	finalOpt := Options{Actions: actions, BonusActions: bonusActions, Passives: passives}

	return finalOpt
}

func GetCharacter(c echo.Context) error {
	db, err := gorm.Open(sqlite.Open("file.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	finalOpt := getCharacterOptions(db, c.Param("name"))

	return c.Render(http.StatusOK, "character", finalOpt)
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
		Name:  hero.Name,
		Class: hero.Class,
		Race:  hero.Race,
		Feats: hero.Feats,
		Items: hero.Items,
	})

	finalOpt := getCharacterOptions(db, hero.Name)

	return c.Render(http.StatusOK, "character", finalOpt)
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
	server.GET("/character/:name", GetCharacter)

	// Start server
	server.Logger.Fatal(server.Start(":1323"))
}
