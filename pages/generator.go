package pages

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/yuriizinets/oxigen/api"
	"go.kyoto.codes/v3/component"
	"go.kyoto.codes/zen/v3/errorsx"
	"go.kyoto.codes/zen/v3/httpx"
	"go.kyoto.codes/zen/v3/slice"
)

type GeneratorState struct {
	common

	Query api.GenerateQuery // Generation query
	Image string            // Resulting image generation url
	Fonts []string          // List of available fonts
}

func Generator(ctx *component.Context) component.State {
	// Initialize state and rendering options
	state := &GeneratorState{}
	// Unpack generation query
	errorsx.Must(0, httpx.Query(ctx.Request.URL.Query()).Unmarshal(&state.Query))
	state.Query.BackgroundFrame = ctx.Request.URL.Query().Get("background.overlay") == "on"
	// Compose generation url
	state.Image = fmt.Sprintf(
		`/api/ogen?%s`,
		ctx.Request.URL.Query().Encode(),
	)
	// Get available fonts
	state.Fonts = slice.Filter(
		slice.Map(errorsx.Must(os.ReadDir("assets/dist/fonts")), func(e fs.DirEntry) string { return e.Name() }), // Read fonts dir, convert entries to strings with names
		func(name string) bool { return strings.HasSuffix(name, ".ttf") },                                        // Filter out fonts (double check)
	)
	// Return
	return state
}
