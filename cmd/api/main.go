package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/prosales/go-api-movies/pkg/i18n"
	"github.com/prosales/go-api-movies/pkg/models"
)

func main() {
	// Definir flags para la línea de comandos
	addr := flag.String("addr", ":8080", "Dirección HTTP")
	apiKey := flag.String("apikey", "", "API Key para OMDB")
	templateDir := flag.String("templates", "./templates", "Ruta a las plantillas")
	staticDir := flag.String("static", "./static", "Ruta a los archivos estáticos")
	localesDir := flag.String("locales", "./locales", "Ruta a los archivos de traducción")
	defaultLang := flag.String("lang", "es", "Idioma predeterminado (es, en)")
	flag.Parse()

	// Verificar que se proporcionó una API key
	if *apiKey == "" {
		// Intentar obtener la API key de una variable de entorno
		*apiKey = os.Getenv("OMDB_API_KEY")
		if *apiKey == "" {
			log.Fatal("OMDB API key no proporcionada. Use --apikey o la variable de entorno OMDB_API_KEY")
		}
	}

	// Inicializar el modelo de películas
	movieModel := models.NewMovieModel(*apiKey)

	// Cargar plantillas
	templates, err := loadTemplates(*templateDir)
	if err != nil {
		log.Fatal(err)
	}

	// Inicializar el traductor
	translator, err := i18n.NewTranslator(*localesDir, *defaultLang)
	if err != nil {
		log.Fatalf("Error al inicializar el traductor: %v", err)
	}

	// Inicializar la aplicación
	app := &application{
		movieModel:  movieModel,
		templates:   templates,
		translator:  translator,
		defaultLang: *defaultLang,
	}

	// Configurar el gestor de archivos estáticos
	fileServer := http.FileServer(http.Dir(*staticDir))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Configurar las rutas
	http.HandleFunc("/", app.homeHandler)
	http.HandleFunc("/search", app.searchHandler)
	http.HandleFunc("/movie", app.movieHandler)
	http.HandleFunc("/change-lang", app.changeLangHandler)

	// Iniciar el servidor
	log.Printf("Iniciando servidor en %s", *addr)
	log.Printf("Usando API key: %s", hideApiKey(*apiKey))
	log.Printf("Directorio de plantillas: %s", filepath.Clean(*templateDir))
	log.Printf("Directorio de archivos estáticos: %s", filepath.Clean(*staticDir))
	log.Printf("Directorio de traducciones: %s", filepath.Clean(*localesDir))
	log.Printf("Idioma predeterminado: %s", *defaultLang)

	err = http.ListenAndServe(*addr, nil)
	log.Fatal(err)
}

// hideApiKey oculta parte de la API key para imprimirla en los logs
func hideApiKey(key string) string {
	if len(key) <= 8 {
		return "*****"
	}
	
	visible := 4
	return key[:visible] + "..." + key[len(key)-visible:]
} 