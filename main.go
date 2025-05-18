package main

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

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

var pythonProcess *exec.Cmd

func startPythonAPI() error {
	// Comando para iniciar el servidor Flask
	pythonProcess = exec.Command("py", "-3.12", "api.py")
	pythonProcess.Stdout = os.Stdout
	pythonProcess.Stderr = os.Stderr
	// Iniciar el proceso en segundo plano
	err := pythonProcess.Start()
	if err != nil {
		return err
	}
	// Esperar un momento para que el servidor Flask se inicie
	time.Sleep(2 * time.Second)
	return nil
}

func stopPythonAPI() error {
	if pythonProcess != nil && pythonProcess.Process != nil {
		// Enviar señal SIGTERM para cerrar el proceso
		err := pythonProcess.Process.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("No python process started")
}

func main() {
	e := echo.New()
	e.Static("/assets", "assets")
	e.Renderer = newTemplate()
	err := startPythonAPI()
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer stopPythonAPI()

	e.GET("/", func(c echo.Context) error {
		e.Logger.Printf("(go): GET /")
		e.Logger.Printf("(py): GET /predict")
		resp, err := http.Get("http://localhost:5000/predict_last")
		if err != nil {
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "error getting python thingy",
			}
			e.Logger.Error(response["Message"], err)
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		defer resp.Body.Close()
		e.Logger.Printf("response: %v", resp)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "error getting response body",
			}
			e.Logger.Error(response["Message"], err)
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		e.Logger.Printf("body: %v", body)
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/",
			"Body":         body,
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
