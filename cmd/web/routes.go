package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Serves static files from the ui/static directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler
	// for all URL paths that start with "/static/". For matching paths, we strip
	// the "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the mux.HandleFunc() method.
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return app.logRequest(secureHeaders(mux))
}
