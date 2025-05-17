package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/prosales/go-api-movies/pkg/i18n"
	"github.com/prosales/go-api-movies/pkg/models"
)

// MockOMDBHandler es un manejador HTTP que simula la API de OMDB
func MockOMDBHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.URL.Query().Get("t") {
	case "Star Wars":
		w.Write([]byte(`{
			"Title": "Star Wars",
			"Year": "1977",
			"Rated": "PG",
			"Released": "25 May 1977",
			"Runtime": "121 min",
			"Genre": "Action, Adventure, Fantasy",
			"Director": "George Lucas",
			"Writer": "George Lucas",
			"Actors": "Mark Hamill, Harrison Ford, Carrie Fisher",
			"Plot": "Luke Skywalker joins forces with a Jedi Knight, a cocky pilot, a Wookiee and two droids to save the galaxy from the Empire's world-destroying battle station, while also attempting to rescue Princess Leia from the mysterious Darth Vader.",
			"Poster": "https://m.media-amazon.com/images/M/MV5BNzVlY2MwMjktM2E4OS00Y2Y3LWE3ZjctYzhkZGM3YzA1ZWM2XkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg",
			"imdbID": "tt0076759",
			"Type": "movie",
			"Response": "True"
		}`))
		return
	}
	
	switch r.URL.Query().Get("s") {
	case "Star":
		w.Write([]byte(`{
			"Search": [
				{
					"Title": "Star Wars",
					"Year": "1977",
					"imdbID": "tt0076759",
					"Type": "movie",
					"Poster": "https://m.media-amazon.com/images/M/MV5BNzVlY2MwMjktM2E4OS00Y2Y3LWE3ZjctYzhkZGM3YzA1ZWM2XkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg"
				},
				{
					"Title": "Star Trek",
					"Year": "2009",
					"imdbID": "tt0796366",
					"Type": "movie",
					"Poster": "https://m.media-amazon.com/images/M/MV5BMjE5NDQ5OTE4Ml5BMl5BanBnXkFtZTcwOTE3NDIzMw@@._V1_SX300.jpg"
				}
			],
			"totalResults": "2",
			"Response": "True"
		}`))
		return
	}
	
	// Por defecto, devolver una respuesta de "no encontrado"
	w.Write([]byte(`{"Response":"False","Error":"Movie not found!"}`))
}

func TestIntegration(t *testing.T) {
	// Crear un servidor mock para OMDB
	mockServer := httptest.NewServer(http.HandlerFunc(MockOMDBHandler))
	defer mockServer.Close()
	
	// Configurar la aplicación para los tests
	t.Run("Application_Setup", func(t *testing.T) {
		// Configurar el traductor
		translator, err := i18n.NewTranslator("../../locales", "es")
		if err != nil {
			t.Fatalf("Error al configurar el traductor: %v", err)
		}
		
		// Configurar el modelo de películas (con la URL del servidor mock)
		movieModel := models.NewMovieModel("test_api_key")
		
		// Cargar las plantillas
		templates, err := loadTemplates("../../templates")
		if err != nil {
			t.Fatalf("Error al cargar las plantillas: %v", err)
		}
		
		// Crear la aplicación
		app := &application{
			movieModel:  movieModel,
			templates:   templates,
			translator:  translator,
			defaultLang: "es",
		}
		
		// Configurar los manejadores HTTP
		mux := http.NewServeMux()
		mux.HandleFunc("/", app.homeHandler)
		mux.HandleFunc("/search", app.searchHandler)
		mux.HandleFunc("/movie", app.movieHandler)
		mux.HandleFunc("/change-lang", app.changeLangHandler)
		
		// Crear un servidor de prueba
		testServer := httptest.NewServer(mux)
		defer testServer.Close()
		
		// Probar la página principal
		resp, err := http.Get(testServer.URL)
		if err != nil {
			t.Fatalf("Error al acceder a la página principal: %v", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("La página principal devolvió un código de estado incorrecto: %d", resp.StatusCode)
		}
	})
}

// TestLoadTemplates prueba la función loadTemplates
func TestLoadTemplates(t *testing.T) {
	templates, err := loadTemplates("../../templates")
	if err != nil {
		t.Fatalf("Error al cargar las plantillas: %v", err)
	}
	
	// Verificar que se cargaron las plantillas esperadas
	expectedTemplates := []string{"home.html", "search.html", "movie.html"}
	for _, name := range expectedTemplates {
		if _, ok := templates[name]; !ok {
			t.Errorf("No se cargó la plantilla %s", name)
		}
	}
} 