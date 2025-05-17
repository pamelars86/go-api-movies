# go-api-movies

Aplicación web en Go para buscar películas utilizando la API de OMDB.

## Características

- Búsqueda de películas por título
- Visualización de detalles de películas
- Interfaz responsive usando Bootstrap
- Caché de resultados para mejorar el rendimiento

## Requisitos

- Go 1.16 o superior
- API Key de OMDB (puedes obtener una en [https://www.omdbapi.com/](https://www.omdbapi.com/))

## Instalación

1. Clona el repositorio:
   ```
   git clone https://github.com/prosales/go-api-movies.git
   cd go-api-movies
   ```

2. Instala las dependencias:
   ```
   make deps
   ```

3. Inicializa la estructura de directorios:
   ```
   make init
   ```

## Uso

### Ejecutar en modo desarrollo

```
make dev API_KEY=tu_api_key
```

o

```
export OMDB_API_KEY=tu_api_key
make dev
```

### Compilar y ejecutar

```
make run API_KEY=tu_api_key
```

### Opciones adicionales

Puedes ver todas las opciones disponibles con:

```
make help
```

## Estructura del proyecto

```
go-api-movies/
├── cmd/
│   └── api/           # Punto de entrada de la aplicación
├── pkg/
│   ├── models/        # Modelos de datos
│   └── omdb/          # Cliente para la API de OMDB
├── static/
│   ├── css/           # Hojas de estilo
│   ├── js/            # JavaScript
│   └── img/           # Imágenes
├── templates/         # Plantillas HTML
├── Makefile           # Comandos de utilidad
└── README.md          # Este archivo
```

## Endpoints

- `GET /` - Página principal
- `GET /search?query=texto` - Búsqueda de películas
- `GET /movie?id=imdbID` - Detalles de una película por ID de IMDB
- `GET /movie?t=título` - Detalles de una película por título

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - vea el archivo LICENSE para más detalles.

## Pruebas

La aplicación incluye pruebas unitarias y de integración para los principales componentes:

- Cliente de la API de OMDB (`pkg/omdb/client_test.go`)
- Modelo de películas (`pkg/models/movies_test.go`)
- Controladores HTTP (`cmd/api/handlers_test.go`)
- Pruebas de integración (`cmd/api/main_test.go`)

Para ejecutar todas las pruebas:

```
make test
```

Para ejecutar pruebas con información de cobertura:

```
make test-cover
```

Para generar un informe HTML de cobertura:

```
make test-cover-report
```

