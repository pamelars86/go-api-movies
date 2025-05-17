package main

import (
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/prosales/go-api-movies/pkg/i18n"
	"github.com/prosales/go-api-movies/pkg/models"
	"github.com/prosales/go-api-movies/pkg/omdb"
)

// Estructura para almacenar el contexto de los handlers
type application struct {
	movieModel  models.MovieModelInterface
	templates   map[string]*template.Template
	translator  i18n.TranslatorInterface
	defaultLang string
}

// Función para renderizar plantillas
func (app *application) render(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	t, ok := app.templates[tmpl]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	// Obtener el idioma desde la cookie o parámetro
	lang := app.getLangFromRequest(r)
	
	// Crear un mapa de funciones para las plantillas que incluya t
	funcMap := template.FuncMap{
		"t": func(key string) string {
			return app.translator.T(lang, key)
		},
	}

	// Si data es un viewData, actualiza el lang
	if vd, ok := data.(*viewData); ok {
		vd.Lang = lang
	}

	// Copiar el template y añadir las funciones de i18n
	// Como no podemos modificar directamente el template, creamos un nuevo
	// buffer para guardar el resultado
	var buf bytes.Buffer
	// Ejecutamos el template con los datos y funciones
	if err := t.Funcs(funcMap).Execute(&buf, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Escribimos el resultado al ResponseWriter
	buf.WriteTo(w)
}

// getLangFromRequest obtiene el idioma de la solicitud (cookie, query param, o header Accept-Language)
func (app *application) getLangFromRequest(r *http.Request) string {
	// 1. Verificar el parámetro de consulta lang
	if lang := r.URL.Query().Get("lang"); lang != "" {
		// No podemos establecer cookies aquí, solo devolver el idioma
		return lang
	}

	// 2. Verificar cookie
	if cookie, err := r.Cookie("lang"); err == nil {
		return cookie.Value
	}

	// 3. Verificar encabezado Accept-Language
	acceptLang := r.Header.Get("Accept-Language")
	if acceptLang != "" {
		// Parsear y obtener el primer idioma preferido
		langs := strings.Split(acceptLang, ",")
		if len(langs) > 0 {
			// Extraer solo el código de idioma (es-ES -> es)
			langCode := strings.Split(langs[0], "-")[0]
			if langCode == "es" || langCode == "en" {
				return langCode
			}
		}
	}

	// 4. Usar idioma predeterminado
	return app.defaultLang
}

// Estructura para datos comunes en las vistas
type viewData struct {
	Error  string
	Query  string
	Movies []omdb.Movie
	Lang   string
	*omdb.Movie
	FromCache  bool
	CachedAt   time.Time
	CacheHits  int
	CacheMisses int
}

// Handler para la página principal
func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := &viewData{
		Lang: app.getLangFromRequest(r),
	}
	app.render(w, r, "home.html", data)
}

// Handler para cambiar el idioma
func (app *application) changeLangHandler(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang != "es" && lang != "en" {
		lang = app.defaultLang
	}

	// Establecer cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "lang",
		Value:  lang,
		Path:   "/",
		MaxAge: 31536000, // Un año
	})

	// Redirigir a la página anterior o a la principal
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

// Handler para la búsqueda
func (app *application) searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	lang := app.getLangFromRequest(r)
	
	data := &viewData{
		Query: query,
		Lang:  lang,
	}

	if query == "" {
		app.render(w, r, "search.html", data)
		return
	}

	result, err := app.movieModel.Search(query)
	if err != nil {
		data.Error = app.translator.T(lang, "error_search") + ": " + err.Error()
		app.render(w, r, "search.html", data)
		return
	}

	if result.Response == "False" || len(result.Search) == 0 {
		data.Movies = []omdb.Movie{}
	} else {
		data.Movies = result.Search
	}

	app.render(w, r, "search.html", data)
}

// Handler para los detalles de una película
func (app *application) movieHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	title := r.URL.Query().Get("t")
	lang := app.getLangFromRequest(r)

	var cachedMovie *models.CachedMovie
	var err error

	// Si se proporciona un ID, redirigimos a la API directamente
	if id != "" {
		http.Redirect(w, r, "https://www.imdb.com/title/"+id, http.StatusSeeOther)
		return
	}

	// Si no hay ID pero hay título, buscamos por título
	if title != "" {
		cachedMovie, err = app.movieModel.GetByTitle(title)
		if err != nil {
			app.render(w, r, "movie.html", &viewData{
				Error: app.translator.T(lang, "error_movie") + ": " + err.Error(),
				Lang:  lang,
			})
			return
		}
	} else {
		// Si no hay ni ID ni título, mostramos un error
		app.render(w, r, "movie.html", &viewData{
			Error: app.translator.T(lang, "error_require_id_title"),
			Lang:  lang,
		})
		return
	}

	// Obtenemos las estadísticas de caché
	hits, misses := app.movieModel.GetCacheStats()

	data := &viewData{
		Movie:      cachedMovie.Movie,
		Lang:       lang,
		FromCache:  cachedMovie.FromCache,
		CachedAt:   cachedMovie.CachedAt,
		CacheHits:  hits,
		CacheMisses: misses,
	}
	app.render(w, r, "movie.html", data)
}

// Función para cargar las plantillas
func loadTemplates(dir string) (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}
	
	// Crear el conjunto de funciones que estarán disponibles para las plantillas
	// Nota: la función real "t" se añade en el momento de renderizar
	funcMap := template.FuncMap{
		"t": func(key string) string {
			return key // Placeholder, será reemplazado en cada renderizado
		},
	}

	pages, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		
		// Saltar layout.html
		if name == "layout.html" {
			continue
		}

		// Crear un nuevo template con las funciones
		ts, err := template.New("layout.html").Funcs(funcMap).ParseFiles(
			filepath.Join(dir, "layout.html"),
			page,
		)
		if err != nil {
			return nil, err
		}

		templates[name] = ts
	}

	return templates, nil
} 