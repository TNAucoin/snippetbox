package main

import (
	"github.com/tnaucoin/snippetbox/internal/models"
	"html/template"
	"path/filepath"
	"time"
)

// Define a templateData type to act as the holding structure for any dynamic data
// that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
	Form        any
}

// humanDate function returns a nicely formatted string representation of a time.Time
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Define a template.FuncMap object which holds the application's dynamic functions.
var functions = template.FuncMap{
	"humanDate": humanDate,
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

		// Parse the base template file, and add any template functions.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
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
