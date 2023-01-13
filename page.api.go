package main

import "git.sr.ht/~kyoto-framework/kyoto"

// PAPIState is a state of PAPI.
type PAPIState struct {
	Routes []PAPIRoute
}

type PAPIRoute struct {
	Method string
	Path   string
	Help   string
	Query  []PAPIQuery
}

type PAPIQuery struct {
	Name     string
	Type     string
	Help     string
	Default  string
	Required bool
}

// PIndex is an API documentation page.
func PAPI(ctx *kyoto.Context) (state PAPIState) {
	// Setup rendering
	kyoto.Template(ctx, "page.api.go.html")
	// Setup API routes doc
	state.Routes = []PAPIRoute{
		{
			Method: "GET",
			Path:   "/api/ogen",
			Help:   "generate social media image",
			Query: []PAPIQuery{
				{
					Name:    "width",
					Type:    "int",
					Help:    "image width",
					Default: "1200",
				},
				{
					Name:    "height",
					Type:    "int",
					Help:    "image height",
					Default: "628",
				},
				{
					Name: "title",
					Type: "string",
					Help: "main image text (left-top aligned, large font)",
				},
				{
					Name: "author",
					Type: "string",
					Help: "source author (left-bottom aligned, medium bold font)",
				},
				{
					Name: "website",
					Type: "string",
					Help: "source link (left-bottom aligned, medium light font)",
				},
				{
					Name: "logo",
					Type: "string",
					Help: "url to source logo (right-bottom aligned, cropped to 250 max width)",
				},
				{
					Name: "background",
					Type: "string",
					Help: "url to image background (black background by default)",
				},
				{
					Name:    "overlay",
					Type:    "bool",
					Help:    "use dark frame on top of the background",
					Default: "false",
				},
				{
					Name:    "dim",
					Type:    "int",
					Help:    "dim background (0-255)",
					Default: "0",
				},
			},
		},
	}
	// Return
	return
}
