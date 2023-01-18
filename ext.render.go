package main

import (
	"image"
	"image/color"
	"io"
	"os"
	"path"

	"github.com/fogleman/gg"
	"github.com/kyoto-framework/zen/v2"
)

// render is an instance of RenderExtension.
// See RenderExtension for details.
var render = &RenderExtension{
	Fonts: "dist/fonts",
}

type Text struct {
	Text  string
	Font  string
	Size  float64
	Align gg.Align
	Color color.Color

	Width float64
}

type Point struct {
	X float64
	Y float64
}

// RenderExtension is a set of helpers
// to reduce complexity of interaction with
// images, text and gg.Context.
type RenderExtension struct {
	Fonts string
}

// LoadRemoteImage downloads remote image,
// saves it into temporary file,
// loads into image object
// and returns that object
// with cleanup function and error.
func (r *RenderExtension) LoadRemoteImage(href string) (image.Image, func(), error) {
	// Fetch image
	res := zen.Request("GET", href).Do()
	// Decode into bytes
	resbts, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	res.Body.Close()
	// Create temp file
	resfile, err := os.CreateTemp("/tmp", "*.oxigen.tmp")
	if err != nil {
		return nil, nil, err
	}
	// Initialize cleanup function
	cleanup := func() {
		os.Remove(resfile.Name())
	}
	// Write bytes into temp file
	if _, err := resfile.Write(resbts); err != nil {
		return nil, nil, err
	}
	// Close file
	if err := resfile.Close(); err != nil {
		return nil, nil, err
	}
	// Load into object
	obj, err := gg.LoadImage(resfile.Name())
	if err != nil {
		return nil, nil, err
	}
	// Return
	return obj, cleanup, nil
}

// Text is a method to render a given text
// with parameters on given context and location.
// Also, provides reasonable defaults for this particular project.
func (r *RenderExtension) Text(img *gg.Context, location Point, text Text) {
	// Defaults
	text.Font = zen.Or(text.Font, "OpenSans-Regular.ttf")
	text.Size = zen.Or(text.Size, 50)
	if zen.Sum(text.Color.RGBA()) == 0 { // If everything is zero, it means no color provided
		text.Color = color.Black
	}
	// Load font
	if err := img.LoadFontFace(path.Join(r.Fonts, text.Font), text.Size); err != nil {
		panic("error while reading font file")
	}
	// Set color
	img.SetColor(text.Color)
	// Draw
	img.DrawStringWrapped(
		text.Text,
		location.X,
		location.Y,
		0, 0,
		text.Width,
		1.5,
		text.Align,
	)
}
