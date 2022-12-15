package main

import "git.sr.ht/~kyoto-framework/kyoto"

type PIndexState struct{}

func PIndex(ctx *kyoto.Context) (state PIndexState) {
	// Setup rendering
	kyoto.Template(ctx, "page.index.go.html")
	// Return
	return
}
