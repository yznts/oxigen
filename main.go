package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"git.sr.ht/~kyoto-framework/kyoto"
	"git.sr.ht/~kyoto-framework/zen"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// setupMiddlewares installs common project middlewares into provided mux.
func setupMiddlewares(mux *mux.Router) {
	mux.Use(func(handler http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, handler)
	})
}

// setupStatic registers a static files handler.
func setupStatic(mux *mux.Router) {
	mux.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./dist"))),
	)
}

// setupKyoto provides advanced configuration for kyoto.
func setupKyoto(mux *mux.Router) {
	kyoto.TemplateConf.FuncMap = kyoto.ComposeFuncMap(
		kyoto.FuncMap, zen.FuncMap,
	)
}

// setupPages registers project pages.
func setupPages(mux *mux.Router) {
	// We are using custom pages register function here.
	// Check Page description for details.
	Page(mux, "/", PIndex)
}

// setupActions registers actions for dynamic components.
func setupActions(mux *mux.Router) {
	// We are using custom actions register function here.
	// Check Action description for details.
	// Action(mux, CExample(nil))
}

func setupAPI(mux *mux.Router) {
	mux.HandleFunc("/api/ogen", AGenerate)
}

// main is a project entry point.
func main() {
	// Parse arguments
	addr := flag.String("http", ":8000", "Serving address")
	flag.Parse()

	// Initialize mux
	mux := mux.NewRouter()

	// Setup parts into mux
	setupMiddlewares(mux)
	setupStatic(mux)
	setupKyoto(mux)
	setupPages(mux)
	setupActions(mux)
	setupAPI(mux)

	// Handle mux into root
	http.Handle("/", mux)

	// Serve
	os.Stdout.WriteString(fmt.Sprintf("Serving on :%s\n", *addr))
	zen.Must(0, http.ListenAndServe(*addr, mux))
}
