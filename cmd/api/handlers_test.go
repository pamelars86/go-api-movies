package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prosales/go-api-movies/pkg/i18n"
	"github.com/prosales/go-api-movies/pkg/models"
	"github.com/prosales/go-api-movies/pkg/omdb"
)

// MockMovieModel es una implementación mock del modelo de películas para pruebas
type MockMovieModel struct {
	GetByTitleFunc   func(title string) (*models.CachedMovie, error)
	SearchFunc       func(query string) (*omdb.SearchResult, error)
	GetCacheStatsFunc func() (hits, misses int)
}

func (m *MockMovieModel) GetByTitle(title string) (*models.CachedMovie, error) {
	return m.GetByTitleFunc(title)
}

func (m *MockMovieModel) Search(query string) (*omdb.SearchResult, error) {
	return m.SearchFunc(query)
}

func (m *MockMovieModel) GetCacheStats() (hits, misses int) {
	return m.GetCacheStatsFunc()
}

// MockTranslator es una implementación mock del traductor para pruebas
type MockTranslator struct {
	TFunc func(lang, key string) string
}

func (m *MockTranslator) T(lang, key string) string {
	return m.TFunc(lang, key)
}

// Test para homeHandler
func TestHomeHandler(t *testing.T) {
	// Crear un traductor mock
	mockTranslator := &MockTranslator{
		TFunc: func(lang, key string) string {
			return key // Devuelve la clave como valor
		},
	}
	
	// Crear plantillas mock
	tmpl, err := template.New("home.html").Parse("{{.Lang}}")
	if err != nil {
		t.Fatal(err)
	}
	templates := map[string]*template.Template{
		"home.html": tmpl,
	}
	
	// Crear la aplicación con el modelo mock
	app := &application{
		translator:  mockTranslator,
		templates:   templates,
		defaultLang: "es",
	}
	
	// Crear una solicitud HTTP mock
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	
	// Llamar al handler
	app.homeHandler(w, req)
	
	// Verificar el código de estado
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

// Test para searchHandler con consulta vacía
func TestSearchHandler_EmptyQuery(t *testing.T) {
	// Crear un traductor mock
	mockTranslator := &MockTranslator{
		TFunc: func(lang, key string) string {
			return key // Devuelve la clave como valor
		},
	}
	
	// Crear plantillas mock
	tmpl, err := template.New("search.html").Parse("{{.Query}}")
	if err != nil {
		t.Fatal(err)
	}
	templates := map[string]*template.Template{
		"search.html": tmpl,
	}
	
	// Crear la aplicación con el modelo mock
	app := &application{
		translator:  mockTranslator,
		templates:   templates,
		defaultLang: "es",
	}
	
	// Crear una solicitud HTTP mock
	req := httptest.NewRequest("GET", "/search", nil)
	w := httptest.NewRecorder()
	
	// Llamar al handler
	app.searchHandler(w, req)
	
	// Verificar el código de estado
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

// Test para searchHandler con consulta
func TestSearchHandler_WithQuery(t *testing.T) {
	// Crear un modelo mock
	mockModel := &MockMovieModel{
		SearchFunc: func(query string) (*omdb.SearchResult, error) {
			// Simular resultados de búsqueda
			return &omdb.SearchResult{
				Search: []omdb.Movie{
					{
						Title:  "Test Movie",
						Year:   "2023",
						ImdbID: "tt1234567",
					},
				},
				TotalResults: "1",
				Response:     "True",
			}, nil
		},
		GetCacheStatsFunc: func() (hits, misses int) {
			return 0, 0
		},
	}
	
	// Crear un traductor mock
	mockTranslator := &MockTranslator{
		TFunc: func(lang, key string) string {
			return key // Devuelve la clave como valor
		},
	}
	
	// Crear plantillas mock
	tmpl, err := template.New("search.html").Parse("{{.Query}}")
	if err != nil {
		t.Fatal(err)
	}
	templates := map[string]*template.Template{
		"search.html": tmpl,
	}
	
	// Crear la aplicación con el modelo mock
	app := &application{
		movieModel: mockModel,
		translator: mockTranslator,
		templates:  templates,
		defaultLang: "es",
	}
	
	// Crear una solicitud HTTP mock
	req := httptest.NewRequest("GET", "/search?query=test", nil)
	w := httptest.NewRecorder()
	
	// Llamar al handler
	app.searchHandler(w, req)
	
	// Verificar el código de estado
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

// Test para movieHandler con título
func TestMovieHandler_WithTitle(t *testing.T) {
	// Crear un modelo mock
	mockModel := &MockMovieModel{
		GetByTitleFunc: func(title string) (*models.CachedMovie, error) {
			// Simular una película
			return &models.CachedMovie{
				Movie: &omdb.Movie{
					Title:    "Test Movie",
					Year:     "2023",
					Director: "Test Director",
					Plot:     "Test Plot",
				},
				FromCache: true,
				CachedAt:  time.Now(),
			}, nil
		},
		GetCacheStatsFunc: func() (hits, misses int) {
			return 1, 0
		},
	}
	
	// Crear un traductor mock
	mockTranslator := &MockTranslator{
		TFunc: func(lang, key string) string {
			return key // Devuelve la clave como valor
		},
	}
	
	// Crear plantillas mock
	tmpl, err := template.New("movie.html").Parse("{{.Movie.Title}}")
	if err != nil {
		t.Fatal(err)
	}
	templates := map[string]*template.Template{
		"movie.html": tmpl,
	}
	
	// Crear la aplicación con el modelo mock
	app := &application{
		movieModel: mockModel,
		translator: mockTranslator,
		templates:  templates,
		defaultLang: "es",
	}
	
	// Crear una solicitud HTTP mock
	req := httptest.NewRequest("GET", "/movie?t=test", nil)
	w := httptest.NewRecorder()
	
	// Llamar al handler
	app.movieHandler(w, req)
	
	// Verificar el código de estado
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

// Test para movieHandler sin título ni ID
func TestMovieHandler_NoTitleOrID(t *testing.T) {
	// Crear un traductor mock
	mockTranslator := &MockTranslator{
		TFunc: func(lang, key string) string {
			return key // Devuelve la clave como valor
		},
	}
	
	// Crear plantillas mock
	tmpl, err := template.New("movie.html").Parse("{{.Error}}")
	if err != nil {
		t.Fatal(err)
	}
	templates := map[string]*template.Template{
		"movie.html": tmpl,
	}
	
	// Crear la aplicación con el modelo mock
	app := &application{
		translator: mockTranslator,
		templates:  templates,
		defaultLang: "es",
	}
	
	// Crear una solicitud HTTP mock
	req := httptest.NewRequest("GET", "/movie", nil)
	w := httptest.NewRecorder()
	
	// Llamar al handler
	app.movieHandler(w, req)
	
	// Verificar el código de estado
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

// Test para changeLangHandler
func TestChangeLangHandler(t *testing.T) {
	// Crear la aplicación
	app := &application{
		defaultLang: "es",
	}
	
	// Crear una solicitud HTTP mock
	req := httptest.NewRequest("GET", "/change-lang?lang=en", nil)
	req.Header.Set("Referer", "/")
	w := httptest.NewRecorder()
	
	// Llamar al handler
	app.changeLangHandler(w, req)
	
	// Verificar el código de estado
	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status code %d, got %d", http.StatusSeeOther, w.Code)
	}
	
	// Verificar que se establece la cookie
	cookies := w.Result().Cookies()
	foundCookie := false
	for _, cookie := range cookies {
		if cookie.Name == "lang" && cookie.Value == "en" {
			foundCookie = true
			break
		}
	}
	if !foundCookie {
		t.Error("Expected lang cookie to be set, but it wasn't")
	}
	
	// Verificar la redirección
	location := w.Header().Get("Location")
	if location != "/" {
		t.Errorf("Expected redirect to /, got %s", location)
	}
}

// Test para getLangFromRequest con cookie
func TestGetLangFromRequest_Cookie(t *testing.T) {
	// Crear la aplicación
	app := &application{
		defaultLang: "es",
	}
	
	// Crear una solicitud HTTP mock con cookie
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "lang",
		Value: "en",
	})
	
	// Obtener el idioma
	lang := app.getLangFromRequest(req)
	
	// Verificar el idioma
	if lang != "en" {
		t.Errorf("Expected lang=en, got %s", lang)
	}
}

// Test para getLangFromRequest con cabecera Accept-Language
func TestGetLangFromRequest_AcceptLanguage(t *testing.T) {
	// Crear la aplicación
	app := &application{
		defaultLang: "es",
	}
	
	// Crear una solicitud HTTP mock con cabecera Accept-Language
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	
	// Obtener el idioma
	lang := app.getLangFromRequest(req)
	
	// Verificar el idioma
	if lang != "en" {
		t.Errorf("Expected lang=en, got %s", lang)
	}
} 