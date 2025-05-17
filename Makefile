.PHONY: build run clean

# Variables
BINARY_NAME=movies-app
API_KEY?=aec2f4c9  # Valor por defecto, reemplazar con tu propia API key

# Construir la aplicación
build:
	@echo "Construyendo la aplicación..."
	go build -o bin/$(BINARY_NAME) ./cmd/api

# Ejecutar la aplicación
run: build
	@echo "Ejecutando la aplicación..."
	./bin/$(BINARY_NAME) --apikey=$(API_KEY)

# Ejecutar en modo desarrollo
dev:
	@echo "Ejecutando en modo desarrollo..."
	go run ./cmd/api --apikey=$(API_KEY)

# Limpiar binarios
clean:
	@echo "Limpiando binarios..."
	rm -rf bin/

# Instalar dependencias
deps:
	@echo "Instalando dependencias..."
	go mod tidy

# Crear la estructura de directorios necesarios
init:
	@echo "Creando estructura de directorios..."
	mkdir -p bin
	mkdir -p static/img
	touch static/img/no-poster.jpg

# Mostrar ayuda
help:
	@echo "Comandos disponibles:"
	@echo "  make build    - Construir la aplicación"
	@echo "  make run      - Construir y ejecutar la aplicación"
	@echo "  make dev      - Ejecutar en modo desarrollo"
	@echo "  make clean    - Limpiar binarios"
	@echo "  make deps     - Instalar dependencias"
	@echo "  make init     - Crear estructura de directorios"
	@echo ""
	@echo "Variables:"
	@echo "  API_KEY       - API Key para OMDB (por defecto: aec2f4c9)"
	@echo ""
	@echo "Ejemplo:"
	@echo "  make run API_KEY=your_api_key"

# Tests
test:
	@echo "Ejecutando pruebas..."
	go test -v ./...

test-cover:
	@echo "Ejecutando pruebas con cobertura..."
	go test -v -cover ./...

test-cover-report:
	@echo "Generando informe de cobertura..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Informe generado en coverage.html" 