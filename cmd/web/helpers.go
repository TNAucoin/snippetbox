package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, td templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}
	// Write the template to a buffer instead of straight to the ResponseWriter
	// so that we can check if there was an error.
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", td)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)
	// Write the contents of the buffer to the http.ResponseWriter
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

// serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError helper sends a specific status code and corresponding description
// to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
