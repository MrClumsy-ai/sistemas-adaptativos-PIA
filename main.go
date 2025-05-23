package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	PORT = ":8080"
	URL  = "http://localhost" + PORT
)

/*
uso en front-end con go:
{{ .Predictions.Apertura.LastDate }}
{{ .Predictions.Apertura.LastValues }}
{{ .Predictions.Apertura.NextPrediction }}

(api.py)
si se quiere conseguir los JSON directamente, usar la url para el get:
http://localhost:5000/

para las predicciones de apertura, usar la ruta:
/predict_last_apertura

para las predicciones de clausura, usar la ruta:
/predict_last_clausura
*/

type Prediction struct {
	LastDate       string    `json:"last_date"`
	LastValues     []float32 `json:"last_values"`
	NextPrediction float32   `json:"next_prediction"`
}

type Predictions struct {
	Apertura Prediction
	Clausura Prediction
}

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
var fetchData *exec.Cmd

func startPythonAPI() error {
	pythonProcess = exec.Command("python3", "api.py")
	pythonProcess.Stdout = os.Stdout
	pythonProcess.Stderr = os.Stderr
	err := pythonProcess.Start()
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	return nil
}

func stopPythonAPI() error {
	if pythonProcess != nil && pythonProcess.Process != nil {
		err := pythonProcess.Process.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("No python process started")
}

func startFetching() error {
	fetchData = exec.Command("python", "database/getDB.py")
	fetchData.Stdout = os.Stdout
	fetchData.Stderr = os.Stderr
	err := fetchData.Start()
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	return nil
}

func stopFetching() error {
	if fetchData != nil && fetchData.Process != nil {
		err := fetchData.Process.Signal(syscall.SIGTERM)
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
	err = startFetching()
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer stopFetching()

	e.GET("/", func(c echo.Context) error {
		fmt.Println("(go): GET /")
		respApertura, err := http.Get("http://localhost:5000/predict_last_apertura")
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "py api: /predict_last_apertura no disponible. Esperar a que se inicie el servidor",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		defer respApertura.Body.Close()
		bodyApertura, err := io.ReadAll(respApertura.Body)
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Error getting apertura body",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		var predictionApertura Prediction
		err = json.Unmarshal(bodyApertura, &predictionApertura)
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Apertura body no pudo ser procesada",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		respClausura, err := http.Get("http://localhost:5000/predict_last_clausura")
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "py api: /predict_last_apertura no disponible. Esperar a que se inicie el servidor",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		defer respClausura.Body.Close()
		bodyClausura, err := io.ReadAll(respClausura.Body)
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Error getting clausura body",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		var predictionClausura Prediction
		err = json.Unmarshal(bodyClausura, &predictionClausura)
		if err != nil {
			e.Logger.Error(err)
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "Apertura body no pudo ser procesada",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}
		predictions := Predictions{
			Apertura: predictionApertura,
			Clausura: predictionClausura,
		}
		response := map[string]any{
			"URL":          URL,
			"CurrentRoute": "/",
			"Predictions":  predictions,
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
