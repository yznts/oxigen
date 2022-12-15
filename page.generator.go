package main

import "git.sr.ht/~kyoto-framework/kyoto"

// PIndexState is a state of PIndex.
type PGeneratorState struct {
	Generator *kyoto.ComponentF[CGeneratorState]
}

// PIndex is a generator UI page.
func PGenerator(ctx *kyoto.Context) (state PGeneratorState) {
	// Setup rendering
	kyoto.Template(ctx, "page.generator.go.html")
	// Init components
	state.Generator = kyoto.Use(ctx, CGenerator(&CGeneratorArgs{}))
	// Return
	return
}
