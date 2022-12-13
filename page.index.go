package main

import "git.sr.ht/~kyoto-framework/kyoto"

// PIndexState is a state of PIndex.
type PIndexState struct {
	Generator *kyoto.ComponentF[CGeneratorState]
}

// PIndex is a home page.
func PIndex(ctx *kyoto.Context) (state PIndexState) {
	// Setup rendering
	kyoto.Template(ctx, "page.index.go.html")
	// Init components
	state.Generator = kyoto.Use(ctx, CGenerator(&CGeneratorArgs{}))
	// Return
	return
}
