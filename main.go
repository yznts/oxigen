package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kyoto-framework/kyoto/v2"
	"github.com/kyoto-framework/zen/v2"
)

// Embed filesystem

//go:embed *.go.html
var fstemplates embed.FS

//go:embed dist
var fsdist embed.FS

// setupMiddlewares installs common project middlewares into provided mux.
func setupMiddlewares(mux *mux.Router) {
	mux.Use(func(handler http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, handler)
	})
}

// setupAssets registers a static files handler.
func setupAssets(mux *mux.Router) {
	switch os.Getenv("DEV") {
	case "true":
		mux.PathPrefix("/assets/").Handler(
			http.StripPrefix("/assets/", http.FileServer(http.Dir("dist"))),
		)
	default:
		distroot, _ := fs.Sub(fsdist, "dist")
		mux.PathPrefix("/assets/").Handler(
			http.StripPrefix("/assets/", http.FileServer(http.FS(distroot))),
		)
	}
}

// setupKyoto provides advanced configuration for kyoto.
func setupKyoto(mux *mux.Router) {
	switch os.Getenv("DEV") {
	case "true":
		kyoto.TemplateConf = kyoto.TemplateConfiguration{
			ParseGlob: "*.go.html",
			FuncMap: kyoto.ComposeFuncMap(
				kyoto.FuncMap, zen.FuncMap, FuncMap,
			),
		}
	default:
		kyoto.TemplateConf = kyoto.TemplateConfiguration{
			ParseFS:   &fstemplates,
			ParseGlob: "*.go.html",
			FuncMap: kyoto.ComposeFuncMap(
				kyoto.FuncMap, zen.FuncMap, FuncMap,
			),
		}
	}
}

// setupPages registers project pages.
func setupPages(mux *mux.Router) {
	// We are using custom pages register function here.
	// Check Page description for details.
	Page(mux, "/", PIndex)
	Page(mux, "/ui", PUI)
	Page(mux, "/api", PAPI)
}

// setupActions registers actions for dynamic components.
func setupActions(mux *mux.Router) {
	// We are using custom actions register function here.
	// Check Action description for details.
	// Action(mux, CExample(nil))
}

// setupAPI registers API handlers.
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
	setupAssets(mux)
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
