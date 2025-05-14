package main

import (
	"database/sql"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

const PORT = "8080"

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	tmpl := template.Must(template.ParseGlob("views/*.tmpl"))
	tmpl = template.Must(tmpl.ParseGlob("views/partials/*.tmpl"))
	return &Templates{
		templates: tmpl,
	}
}

func main() {
	e := echo.New()
	e.Static("/assets", "assets")
	e.Renderer = newTemplate()
	dbConnection, err := sql.Open("sqlite3", "pia.db")
	if err != nil {
		e.Logger.Fatal("error connecting to db: ", err)
	}
	defer dbConnection.Close()
	/* _, err = dbConnection.Exec(`CREATE TABLE IF NOT EXISTS mascotas (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	nombre TEXT NOT NULL,
	edad INTEGER,
	altura_cm INTEGER,
	foto BLOB,
	descripcion TEXT)`)
	if err != nil {
		e.Logger.Fatal("error creating mascota table", err)
	} */

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "inicio", nil)
	})
	e.Logger.Fatal(e.Start(":" + PORT))
}
