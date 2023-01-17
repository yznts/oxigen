package main

import "github.com/kyoto-framework/kyoto/v2"

type PIndexState struct{}

func PIndex(ctx *kyoto.Context) (state PIndexState) {
	// Setup rendering
	kyoto.Template(ctx, "page.index.go.html")
	// Return
	return
}
