package i18n

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// TranslatorInterface define la interfaz para un traductor
type TranslatorInterface interface {
	T(lang, key string) string
}

// Translator maneja las traducciones i18n
type Translator struct {
	translations map[string]map[string]string
	defaultLang  string
	mu           sync.RWMutex
}

// NewTranslator crea un nuevo traductor con el idioma predeterminado
func NewTranslator(localesDir, defaultLang string) (*Translator, error) {
	t := &Translator{
		translations: make(map[string]map[string]string),
		defaultLang:  defaultLang,
	}

	// Cargar todos los idiomas disponibles en el directorio locales
	langDirs, err := ioutil.ReadDir(localesDir)
	if err != nil {
		return nil, err
	}

	for _, langDir := range langDirs {
		if !langDir.IsDir() {
			continue
		}

		lang := langDir.Name()
		// Cargar archivo messages.json para este idioma
		messagesFile := filepath.Join(localesDir, lang, "messages.json")
		if _, err := os.Stat(messagesFile); os.IsNotExist(err) {
			continue
		}

		data, err := ioutil.ReadFile(messagesFile)
		if err != nil {
			return nil, err
		}

		var messages map[string]string
		if err := json.Unmarshal(data, &messages); err != nil {
			return nil, err
		}

		t.translations[lang] = messages
	}

	return t, nil
}

// T traduce una clave al idioma especificado
func (t *Translator) T(lang, key string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Primero intenta encontrar la traducción en el idioma solicitado
	if translations, ok := t.translations[lang]; ok {
		if val, ok := translations[key]; ok {
			return val
		}
	}

	// Si no se encuentra, intenta con el idioma predeterminado
	if translations, ok := t.translations[t.defaultLang]; ok {
		if val, ok := translations[key]; ok {
			return val
		}
	}

	// Si todo falla, devuelve la clave como respaldo
	return key
} 