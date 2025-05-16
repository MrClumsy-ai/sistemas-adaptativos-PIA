package main

import (
	"database/sql"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

const PORT = ":8080"
const URL = "http://localhost" + PORT

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
	const dbLocation = "pia.db"
	dbConnection, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		e.Logger.Fatal("Error connecting to database:", err)
	}
	e.Logger.Printf("Database connection established in %v", dbLocation)
	defer dbConnection.Close()
	/* _, err = dbConnection.Exec(`CREATE TABLE IF NOT EXISTS myTable (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	nombre TEXT NOT NULL,
	edad INTEGER,
	altura_cm INTEGER,
	foto BLOB,
	descripcion TEXT)`)
	if err != nil {
		e.Logger.Fatal("Error creating myTable table", err)
	} */

	e.GET("/", func(c echo.Context) error {
		/* ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		rows, err := dbConnection.QueryContext(ctx, "SELECT * FROM myTable")
		if err != nil {
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Internal server error",
			}
			e.Logger.Error(err)
			return c.Render(http.StatusInternalServerError, "error", response)
		} */
		e.Logger.Printf("GET request for /")
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/",
		}
		return c.Render(http.StatusOK, "inicio", response)
	})

	// catch all route
	e.GET("/*", func(c echo.Context) error {
		response := map[string]any{
			"URL":     URL,
			"Code":    http.StatusNotFound,
			"Message": "Not found",
		}
		return c.Render(http.StatusNotFound, "error", response)
	})
	e.Logger.Fatal(e.Start(PORT))
}
