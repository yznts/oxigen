package pages

import "go.kyoto.codes/v3/component"

type ApiRoute struct {
	Method string
	Path   string
	Help   string
	Query  []ApiQuery
}

type ApiQuery struct {
	Name     string
	Type     string
	Help     string
	Default  string
	Required bool
}

type ApiState struct {
	common

	Routes []ApiRoute
}

func Api(ctx *component.Context) component.State {
	return &ApiState{
		Routes: []ApiRoute{
			{
				Method: "GET",
				Path:   "/api/ogen",
				Help:   "generate social media image",
				Query: []ApiQuery{
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
						Name:    "title.font",
						Type:    "string",
						Default: "OpenSans-SemiBold.ttf",
						Help:    "title font",
					},
					{
						Name:    "title.font.size",
						Type:    "number",
						Default: "80",
						Help:    "title font size",
					},
					{
						Name:    "author.font",
						Type:    "string",
						Default: "OpenSans-SemiBold.ttf",
						Help:    "title font",
					},
					{
						Name:    "author.font.size",
						Type:    "number",
						Default: "50",
						Help:    "author font size",
					},
					{
						Name:    "website.font",
						Type:    "string",
						Default: "OpenSans-Light.ttf",
						Help:    "title font",
					},
					{
						Name:    "website.font.size",
						Type:    "number",
						Default: "50",
						Help:    "author font size",
					},
					{
						Name:    "background.dim",
						Type:    "int",
						Help:    "dim background (0-255)",
						Default: "0",
					},
					{
						Name:    "background.frame",
						Type:    "bool",
						Help:    "use dark frame on top of the background",
						Default: "false",
					},
				},
			},
		},
	}
}
