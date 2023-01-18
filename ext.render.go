package main

import (
	"image"
	"io"
	"os"

	"github.com/fogleman/gg"
	"github.com/kyoto-framework/zen/v2"
)

var extrender = &RenderExtension{}

type RenderExtension struct{}

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
