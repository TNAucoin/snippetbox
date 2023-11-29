package main

import (
	"github.com/tnaucoin/snippetbox/internal/models"
	"html/template"
	"path/filepath"
)

// Define a templateData type to act as the holding structure for any dynamic data
// that we want to pass to our HTML templates.
type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

// newTemplateCache function returns a map of template pointers.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	// Glob the page templates to get a slice of file paths.
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl.html') from the full path.
		name := filepath.Base(page)

		// Parse the base template file
		ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}
		// Parse the partials into the base template
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}
		// Parse the page template into the base template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
