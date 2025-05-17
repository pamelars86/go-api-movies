package omdb

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "test_key"
	client := NewClient(apiKey)

	if client.ApiKey != apiKey {
		t.Errorf("Expected ApiKey to be %s, got %s", apiKey, client.ApiKey)
	}

	if client.HttpClient == nil {
		t.Error("Expected HttpClient to be initialized, got nil")
	}
}

func TestSearchByTitle(t *testing.T) {
	// Crear un servidor de prueba
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar que el método sea correcto
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verificar que los parámetros sean correctos
		q := r.URL.Query()
		if q.Get("apikey") != "test_key" {
			t.Errorf("Expected apikey=test_key, got %s", q.Get("apikey"))
		}
		if q.Get("s") != "test_movie" {
			t.Errorf("Expected s=test_movie, got %s", q.Get("s"))
		}

		// Retornar una respuesta de ejemplo
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"Search": [
				{
					"Title": "Test Movie",
					"Year": "2023",
					"imdbID": "tt1234567",
					"Type": "movie",
					"Poster": "https://example.com/poster.jpg"
				}
			],
			"totalResults": "1",
			"Response": "True"
		}`))
	}))
	defer server.Close()

	// Crear un cliente que use el servidor de prueba
	client := &Client{
		ApiKey:     "test_key",
		HttpClient: server.Client(),
	}

	// Reemplazar la URL base con la del servidor de prueba
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	// Realizar la búsqueda
	result, err := client.SearchByTitle("test_movie")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	// Verificar los resultados
	if result.Response != "True" {
		t.Errorf("Expected Response=True, got %s", result.Response)
	}
	if len(result.Search) != 1 {
		t.Errorf("Expected 1 search result, got %d", len(result.Search))
	}
	if result.Search[0].Title != "Test Movie" {
		t.Errorf("Expected Title=Test Movie, got %s", result.Search[0].Title)
	}
}

func TestGetMovieByTitle(t *testing.T) {
	// Crear un servidor de prueba
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar que el método sea correcto
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verificar que los parámetros sean correctos
		q := r.URL.Query()
		if q.Get("apikey") != "test_key" {
			t.Errorf("Expected apikey=test_key, got %s", q.Get("apikey"))
		}
		if q.Get("t") != "test_movie" {
			t.Errorf("Expected t=test_movie, got %s", q.Get("t"))
		}

		// Retornar una respuesta de ejemplo
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"Title": "Test Movie",
			"Year": "2023",
			"Rated": "PG-13",
			"Released": "01 Jan 2023",
			"Runtime": "120 min",
			"Genre": "Action, Adventure",
			"Director": "Test Director",
			"Writer": "Test Writer",
			"Actors": "Actor 1, Actor 2",
			"Plot": "This is a test movie plot.",
			"Poster": "https://example.com/poster.jpg",
			"imdbID": "tt1234567",
			"Type": "movie",
			"Response": "True"
		}`))
	}))
	defer server.Close()

	// Crear un cliente que use el servidor de prueba
	client := &Client{
		ApiKey:     "test_key",
		HttpClient: server.Client(),
	}

	// Reemplazar la URL base con la del servidor de prueba
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	// Realizar la búsqueda
	movie, err := client.GetMovieByTitle("test_movie")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	// Verificar los resultados
	if movie.Title != "Test Movie" {
		t.Errorf("Expected Title=Test Movie, got %s", movie.Title)
	}
	if movie.Year != "2023" {
		t.Errorf("Expected Year=2023, got %s", movie.Year)
	}
	if movie.Director != "Test Director" {
		t.Errorf("Expected Director=Test Director, got %s", movie.Director)
	}
	if movie.Response != "True" {
		t.Errorf("Expected Response=True, got %s", movie.Response)
	}
}

func TestGetMovieByTitle_Error(t *testing.T) {
	// Crear un servidor de prueba que devuelve un error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retornar una respuesta de error
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"Response": "False",
			"Error": "Movie not found!"
		}`))
	}))
	defer server.Close()

	// Crear un cliente que use el servidor de prueba
	client := &Client{
		ApiKey:     "test_key",
		HttpClient: server.Client(),
	}

	// Reemplazar la URL base con la del servidor de prueba
	originalBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = originalBaseURL }()

	// Realizar la búsqueda
	_, err := client.GetMovieByTitle("nonexistent_movie")

	// Verificar que se retorne un error
	if err == nil {
		t.Error("Expected an error, got nil")
	}
} 