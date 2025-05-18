package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
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

	e.GET("/", func(c echo.Context) error {
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
