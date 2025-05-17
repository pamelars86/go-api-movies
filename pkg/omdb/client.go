package omdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var (
	// BaseURL es la URL base de la API de OMDB (exportable para pruebas)
	BaseURL = "https://www.omdbapi.com/"
)

// Client representa un cliente para la API de OMDB
type Client struct {
	ApiKey     string
	HttpClient *http.Client
}

// Movie representa la estructura de datos de una película de OMDB
type Movie struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Poster   string `json:"Poster"`
	ImdbID   string `json:"imdbID"`
	Type     string `json:"Type"`
	Response string `json:"Response"`
}

// SearchResult representa el resultado de una búsqueda de películas
type SearchResult struct {
	Search       []Movie `json:"Search"`
	TotalResults string  `json:"totalResults"`
	Response     string  `json:"Response"`
}

// OMDBClient define la interfaz para un cliente de OMDB
type OMDBClient interface {
	SearchByTitle(title string) (*SearchResult, error)
	GetMovieByTitle(title string) (*Movie, error)
}

// NewClient crea un nuevo cliente para la API de OMDB
func NewClient(apiKey string) *Client {
	return &Client{
		ApiKey:     apiKey,
		HttpClient: &http.Client{},
	}
}

// SearchByTitle busca películas por título
func (c *Client) SearchByTitle(title string) (*SearchResult, error) {
	params := url.Values{}
	params.Add("apikey", c.ApiKey)
	params.Add("s", title)

	fullURL := fmt.Sprintf("%s?%s", BaseURL, params.Encode())
	
	resp, err := c.HttpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error al hacer la solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code inesperado: %d", resp.StatusCode)
	}

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta: %w", err)
	}

	return &result, nil
}

// GetMovieByTitle obtiene una película por título
func (c *Client) GetMovieByTitle(title string) (*Movie, error) {
	params := url.Values{}
	params.Add("apikey", c.ApiKey)
	params.Add("t", title)

	fullURL := fmt.Sprintf("%s?%s", BaseURL, params.Encode())
	
	resp, err := c.HttpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error al hacer la solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code inesperado: %d", resp.StatusCode)
	}

	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta: %w", err)
	}

	if movie.Response == "False" {
		return nil, fmt.Errorf("película no encontrada")
	}

	return &movie, nil
} 