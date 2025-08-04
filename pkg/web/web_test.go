package web

import (
	"html/template"
	"testing"
)

func TestTemplatesValid(t *testing.T) {
	_, err := template.New("").ParseGlob("../../templates/*.html")
	if err != nil {
		t.Errorf("couldn't initialize templates: %v", err)
	}
}
