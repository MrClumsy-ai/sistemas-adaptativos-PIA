package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	tflite "github.com/mattn/go-tflite"
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

	e.GET("/", func(c echo.Context) error {
		modelPath := "entrega-2/modelo_apertura.tflite"
		csvPath := "entrega-2/corr_bitcoin_diario_apertura.csv"
		e.Logger.Printf("modelo: %v\ncsv: %v", modelPath, csvPath)
		predictions, err := PredictFromCSV(modelPath, csvPath)
		if err != nil {
			response := map[string]any{
				"URL":     URL,
				"Code":    http.StatusInternalServerError,
				"Message": "error csv",
			}
			return c.Render(http.StatusInternalServerError, "error", response)
		}

		// Imprimir resultados
		e.Logger.Printf("\nPredicciones:")
		for i, pred := range predictions {
			fmt.Printf("Fila %d: %v\n", i+1, pred)
		}

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

// PredictFromCSV carga un modelo TFLite y realiza predicciones sobre datos de un CSV
func PredictFromCSV(modelPath string, csvPath string) ([][]float32, error) {
	// 1. Cargar el modelo TensorFlow Lite
	model := tflite.NewModelFromFile(modelPath)
	if model == nil {
		return nil, fmt.Errorf("no se pudo cargar el modelo desde %s", modelPath)
	}
	defer model.Delete()

	// 2. Configurar el intérprete
	options := tflite.NewInterpreterOptions()
	defer options.Delete()
	// options.SetNumThreads(4) // Usar 4 hilos para mejor rendimiento

	interpreter := tflite.NewInterpreter(model, options)
	if interpreter == nil {
		return nil, fmt.Errorf("no se pudo crear el intérprete TFLite")
	}
	defer interpreter.Delete()

	// 3. Asignar memoria para los tensores
	if status := interpreter.AllocateTensors(); status != tflite.OK {
		return nil, fmt.Errorf("fallo al asignar tensores")
	}

	// 4. Obtener detalles del tensor de entrada
	input := interpreter.GetInputTensor(0)
	inputSize := input.ByteSize() / 4 // 4 bytes por float32
	inputShape := input.Shape()
	fmt.Printf("Modelo espera entrada con shape: %v\n", inputShape)

	// 5. Cargar y procesar el archivo CSV
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error al leer CSV: %v", err)
	}

	// 6. Preparar slice para almacenar predicciones
	var predictions [][]float32

	// 7. Procesar cada fila del CSV
	for _, record := range records {
		// Convertir strings a float32
		var inputData []float32
		for _, value := range record {
			num, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return nil, fmt.Errorf("error al convertir valor '%s' a float32: %v", value, err)
			}
			inputData = append(inputData, float32(num))
		}

		// Verificar que los datos coincidan con lo esperado por el modelo
		if len(inputData) != int(inputSize) {
			return nil, fmt.Errorf("tamaño de entrada incorrecto: esperado %d, obtenido %d", inputSize, len(inputData))
		}

		// Copiar datos al tensor de entrada
		copy(input.Float32s(), inputData)

		// Ejecutar inferencia
		if status := interpreter.Invoke(); status != tflite.OK {
			return nil, fmt.Errorf("fallo en la inferencia")
		}

		// Obtener resultados
		output := interpreter.GetOutputTensor(0)
		outputData := output.Float32s()

		// Hacer una copia de los resultados (outputData es un slice que puede cambiar)
		prediction := make([]float32, len(outputData))
		copy(prediction, outputData)

		predictions = append(predictions, prediction)
	}

	return predictions, nil
}
