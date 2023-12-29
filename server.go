package main

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GetIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func GetCharacter(c echo.Context) error {
	db := getDatabaseSession()
	characterInfo := GetCharacterInfo(db, c.QueryParam("name"))

	return c.Render(http.StatusOK, "character", characterInfo)
}

func GetCharacterNames(c echo.Context) error {
	db := getDatabaseSession()

	var names []string
	db.Table("characters").Select("name").Scan(&names)

	return c.Render(http.StatusOK, "drop_down", names)
}

type OptionWithDescription struct {
	Name        string
	Properties  string
	Description string
}

func GetOptionDescription(c echo.Context) error {
	db := getDatabaseSession()

	name := c.Request().Header["Hx-Trigger-Name"][0]

	var description string
	db.Table("options").Select("description").Where("name = ?", name).Scan(&description)

	return c.Render(http.StatusOK, "description", OptionWithDescription{name, "\n", description})
}

func GetItemDescription(c echo.Context) error {
	db := getDatabaseSession()

	name := c.Request().Header["Hx-Trigger-Name"][0]

	var record Item
	db.Table("items").Select("description", "properties").Where("name = ?", name).Scan(&record)

	return c.Render(http.StatusOK, "description", OptionWithDescription{name, record.Properties, record.Description})
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//go:embed web/*
var content embed.FS

func startWebServer() {
	// Echo instance
	server := echo.New()

	absPath, _ := filepath.Abs("web")
	// Middleware
	server.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.RequestID(),
		middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       absPath,
			Filesystem: http.FS(content),
			Browse:     true,
			IgnoreBase: true,
		},
		))

	server.HTTPErrorHandler = func(err error, c echo.Context) {
		// Take required information from error and context and send it to a service like New Relic
		fmt.Println(c.Path(), c.QueryParams(), err.Error())

		// Call the default handler to return the HTTP response
		server.DefaultHTTPErrorHandler(err, c)
	}

	// This will initiate our template renderer
	t := &Template{
		templates: template.Must(template.ParseGlob("web/*.html")),
	}
	server.Renderer = t

	server.GET("/", GetIndex)
	server.File("/web/styles.css", "web/styles.css")
	server.File("/favicon.ico", "web/favicon.ico")

	//// character
	server.GET("/character", GetCharacter)
	server.GET("/character/names", GetCharacterNames)
	server.GET("/option", GetOptionDescription)
	server.GET("/item", GetItemDescription)

	// Start server
	server.Logger.Fatal(server.Start(":1323"))
}
