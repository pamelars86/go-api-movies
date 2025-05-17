package models

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/prosales/go-api-movies/pkg/omdb"
)

// MovieModelInterface define la interfaz para un modelo de películas
type MovieModelInterface interface {
	GetByTitle(title string) (*CachedMovie, error)
	Search(query string) (*omdb.SearchResult, error)
	GetCacheStats() (hits, misses int)
}

// MovieModel representa un modelo para acceder y manipular datos de películas
type MovieModel struct {
	client       omdb.OMDBClient
	cache        map[string]*CachedMovie
	mu           sync.RWMutex
	cacheHits    int
	cacheMisses  int
}

// CachedMovie representa una película con metadatos de caché
type CachedMovie struct {
	Movie      *omdb.Movie
	FromCache  bool
	CachedAt   time.Time
}

// NewMovieModel crea un nuevo modelo de películas
func NewMovieModel(apiKey string) *MovieModel {
	return &MovieModel{
		client: omdb.NewClient(apiKey),
		cache:  make(map[string]*CachedMovie),
	}
}

// GetCacheStats devuelve estadísticas del uso de caché
func (m *MovieModel) GetCacheStats() (hits, misses int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cacheHits, m.cacheMisses
}

// GetByTitle obtiene una película por su título
func (m *MovieModel) GetByTitle(title string) (*CachedMovie, error) {
	if title == "" {
		return nil, errors.New("título vacío")
	}

	// Primero verificamos en la caché
	m.mu.RLock()
	if cachedMovie, ok := m.cache[title]; ok {
		m.cacheHits++
		log.Printf("CACHÉ: Película encontrada en caché: %s", title)
		
		// Marcar como proveniente de caché
		cachedMovie.FromCache = true
		
		m.mu.RUnlock()
		return cachedMovie, nil
	}
	m.mu.RUnlock()

	// Si no está en la caché, lo buscamos en la API
	m.mu.Lock()
	m.cacheMisses++
	m.mu.Unlock()
	
	log.Printf("API: Buscando película en API externa: %s", title)
	movie, err := m.client.GetMovieByTitle(title)
	if err != nil {
		return nil, err
	}

	// Creamos un objeto CachedMovie con metadatos
	cachedMovie := &CachedMovie{
		Movie:     movie,
		FromCache: false,
		CachedAt:  time.Now(),
	}

	// Guardamos en la caché
	m.mu.Lock()
	m.cache[title] = cachedMovie
	m.mu.Unlock()

	return cachedMovie, nil
}

// Search busca películas que coincidan con el término de búsqueda
func (m *MovieModel) Search(query string) (*omdb.SearchResult, error) {
	if query == "" {
		return nil, errors.New("consulta vacía")
	}

	log.Printf("API: Buscando película(s) con consulta: %s", query)
	result, err := m.client.SearchByTitle(query)
	if err != nil {
		return nil, err
	}

	return result, nil
} 