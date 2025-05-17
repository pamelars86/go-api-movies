package models

import (
	"testing"
	"time"

	"github.com/prosales/go-api-movies/pkg/omdb"
)

// MockClient es una implementación mock del cliente OMDB para pruebas
type MockClient struct {
	SearchByTitleFunc   func(title string) (*omdb.SearchResult, error)
	GetMovieByTitleFunc func(title string) (*omdb.Movie, error)
}

func (m *MockClient) SearchByTitle(title string) (*omdb.SearchResult, error) {
	return m.SearchByTitleFunc(title)
}

func (m *MockClient) GetMovieByTitle(title string) (*omdb.Movie, error) {
	return m.GetMovieByTitleFunc(title)
}

// Test para GetByTitle cuando la película está en caché
func TestGetByTitle_FromCache(t *testing.T) {
	// Crear un cliente mock
	mockClient := &MockClient{
		GetMovieByTitleFunc: func(title string) (*omdb.Movie, error) {
			t.Error("No debería llamar a GetMovieByTitle cuando la película está en caché")
			return nil, nil
		},
	}

	// Crear el modelo con el cliente mock
	model := &MovieModel{
		client: mockClient,
		cache:  make(map[string]*CachedMovie),
	}

	// Preparar datos en caché
	testMovie := &omdb.Movie{
		Title: "Cached Movie",
		Year:  "2023",
	}
	model.cache["test_movie"] = &CachedMovie{
		Movie:    testMovie,
		CachedAt: time.Now(),
	}

	// Realizar la búsqueda
	result, err := model.GetByTitle("test_movie")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	// Verificar que se devolvió la película desde la caché
	if result.Movie.Title != "Cached Movie" {
		t.Errorf("Expected Title=Cached Movie, got %s", result.Movie.Title)
	}
	if !result.FromCache {
		t.Error("Expected FromCache=true, got false")
	}
	if model.cacheHits != 1 {
		t.Errorf("Expected cacheHits=1, got %d", model.cacheHits)
	}
}

// Test para GetByTitle cuando la película no está en caché
func TestGetByTitle_FromAPI(t *testing.T) {
	// Crear un cliente mock
	mockClient := &MockClient{
		GetMovieByTitleFunc: func(title string) (*omdb.Movie, error) {
			// Simular una respuesta de la API
			return &omdb.Movie{
				Title: "API Movie",
				Year:  "2023",
			}, nil
		},
	}

	// Crear el modelo con el cliente mock
	model := &MovieModel{
		client: mockClient,
		cache:  make(map[string]*CachedMovie),
	}

	// Realizar la búsqueda
	result, err := model.GetByTitle("new_movie")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	// Verificar que se devolvió la película desde la API
	if result.Movie.Title != "API Movie" {
		t.Errorf("Expected Title=API Movie, got %s", result.Movie.Title)
	}
	if result.FromCache {
		t.Error("Expected FromCache=false, got true")
	}
	if model.cacheMisses != 1 {
		t.Errorf("Expected cacheMisses=1, got %d", model.cacheMisses)
	}

	// Verificar que la película se guardó en caché
	cached, exists := model.cache["new_movie"]
	if !exists {
		t.Error("Expected movie to be cached, but it wasn't")
	}
	if cached.Movie.Title != "API Movie" {
		t.Errorf("Expected cached Title=API Movie, got %s", cached.Movie.Title)
	}
}

// Test para Search
func TestSearch(t *testing.T) {
	// Crear un cliente mock
	mockClient := &MockClient{
		SearchByTitleFunc: func(title string) (*omdb.SearchResult, error) {
			// Simular una respuesta de la API
			return &omdb.SearchResult{
				Search: []omdb.Movie{
					{
						Title:  "Search Result 1",
						Year:   "2021",
						ImdbID: "tt1111111",
					},
					{
						Title:  "Search Result 2",
						Year:   "2022",
						ImdbID: "tt2222222",
					},
				},
				TotalResults: "2",
				Response:     "True",
			}, nil
		},
	}

	// Crear el modelo con el cliente mock
	model := &MovieModel{
		client: mockClient,
		cache:  make(map[string]*CachedMovie),
	}

	// Realizar la búsqueda
	result, err := model.Search("test_search")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	// Verificar los resultados
	if result.Response != "True" {
		t.Errorf("Expected Response=True, got %s", result.Response)
	}
	if len(result.Search) != 2 {
		t.Errorf("Expected 2 search results, got %d", len(result.Search))
	}
	if result.Search[0].Title != "Search Result 1" {
		t.Errorf("Expected Title=Search Result 1, got %s", result.Search[0].Title)
	}
	if result.Search[1].Title != "Search Result 2" {
		t.Errorf("Expected Title=Search Result 2, got %s", result.Search[1].Title)
	}
}

// Test para NewMovieModel
func TestNewMovieModel(t *testing.T) {
	model := NewMovieModel("test_api_key")

	if model.client == nil {
		t.Error("Expected client to be initialized, got nil")
	}
	if model.cache == nil {
		t.Error("Expected cache to be initialized, got nil")
	}
	if model.cacheHits != 0 {
		t.Errorf("Expected cacheHits=0, got %d", model.cacheHits)
	}
	if model.cacheMisses != 0 {
		t.Errorf("Expected cacheMisses=0, got %d", model.cacheMisses)
	}
}

// Test para GetCacheStats
func TestGetCacheStats(t *testing.T) {
	model := &MovieModel{
		client:      &MockClient{},
		cache:       make(map[string]*CachedMovie),
		cacheHits:   5,
		cacheMisses: 3,
	}

	hits, misses := model.GetCacheStats()
	if hits != 5 {
		t.Errorf("Expected hits=5, got %d", hits)
	}
	if misses != 3 {
		t.Errorf("Expected misses=3, got %d", misses)
	}
} 