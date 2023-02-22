package main

import (
	"fmt"

	"github.com/kyoto-framework/kyoto/v2"
	"github.com/kyoto-framework/zen/v3/errorsx"
	"github.com/kyoto-framework/zen/v3/httpx"
)

type CGeneratorArgs struct{}

type CGeneratorState struct {
	Args *CGeneratorArgs

	Query GenerateQuery
	Image string
}

func CGenerator(args *CGeneratorArgs) kyoto.Component[CGeneratorState] {
	return func(ctx *kyoto.Context) (state CGeneratorState) {
		// Write args
		state.Args = args
		// Unpack query
		errorsx.Must(0, httpx.Query(ctx.Request.URL.Query()).Unmarshal(&state.Query))
		// Compose generation url
		state.Image = fmt.Sprintf(
			`/api/ogen?%s`,
			ctx.Request.URL.Query().Encode(),
		)
		// Return
		return
	}
}
