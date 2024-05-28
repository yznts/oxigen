package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/yznts/oxigen/api"
	"github.com/yznts/oxigen/pages"
	"go.kyoto.codes/v3/component"
	"go.kyoto.codes/v3/htmx"
	"go.kyoto.codes/v3/rendering"
	"go.kyoto.codes/zen/v3/errorsx"
	"go.kyoto.codes/zen/v3/mapx"
	"go.kyoto.codes/zen/v3/templatex"
)

func main() {
	// Parse arguments
	addr := flag.String("http", ":8000", "Serving address")
	flag.Parse()

	// Setup rendering
	rendering.TEMPLATE_GLOB = "**/*.go.html"
	rendering.TEMPLATE_FUNCMAP = mapx.Merge(
		rendering.FuncMap,
		htmx.FuncMap,
		component.FuncMap,
		templatex.FuncMap,
	)

	// Initialize mux
	mux := http.NewServeMux()

	// Setup assets
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/dist"))))

	// Setup serving components (pages/components)
	serve := map[string]component.Component{
		"/":          pages.Home,
		"/api":       pages.Api,
		"/generator": pages.Generator,
	}
	// Register serving components
	for route, component := range serve {
		mux.HandleFunc(route, rendering.Handler(component))
	}

	// Setup API
	mux.HandleFunc("/api/ogen", api.Generator)

	// Serve
	log.Printf(fmt.Sprintf("Serving on %s", *addr))
	errorsx.Must(0, http.ListenAndServe(*addr, mux))
}
