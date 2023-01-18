package main

import (
	"image/color"
	"net/http"
	"os"
	"path"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/kyoto-framework/zen/v2"
)

// Generation defaults
var (
	titleFontDefault   = "OpenSans-SemiBold.ttf"
	authorFontDefault  = "OpenSans-SemiBold.ttf"
	websiteFontDefault = "OpenSans-Light.ttf"

	titleFontSizeDefault   = 80.0
	authorFontSizeDefault  = 50.0
	websiteFontSizeDefault = 50.0
)

// Generation constants
var (
	marginOverlay  = 20.0
	marginTitleX   = 60.0
	marginTitleY   = 90.0
	marginWebsiteX = 70.0
	marginWebsiteY = 30.0
	marginAuthorX  = 70.0
	marginAuthorY  = 130.0
	marginLogoX    = 50.0
	marginLogoY    = 40.0

	colorTitle   = color.White
	colorAuthor  = color.White
	colorWebsite = color.White
	colorOverlay = color.RGBA{0, 0, 0, 194}
)

// GenerateQuery holds generation parameters
type GenerateQuery struct {
	// Image parameters
	Width  int `query:"width"`
	Height int `query:"height"`
	// Basic settings
	Title      string `query:"title"`
	Author     string `query:"author"`
	Website    string `query:"website"`
	Logo       string `query:"logo"`
	Background string `query:"background"`
	// Advanced settings
	TitleFont         string  `query:"title.font"`
	TitleFontSize     float64 `query:"title.font.size"`
	AuthorFont        string  `query:"author.font"`
	AuthorFontSize    float64 `query:"author.font.size"`
	WebsiteFont       string  `query:"website.font"`
	WebsiteFontSize   float64 `query:"website.font.size"`
	BackgroundDim     int     `query:"background.dim"`
	BackgroundOverlay bool    `query:"background.overlay"`
}

func AGenerate(w http.ResponseWriter, r *http.Request) {
	// Unpack query
	query := GenerateQuery{}
	zen.Must(0, zen.Query(r.URL.Query()).Unmarshal(&query))
	// Resolve defaults
	query.Width = zen.Or(query.Width, 1200)
	query.Height = zen.Or(query.Height, 628)
	query.TitleFont = zen.Or(query.TitleFont, titleFontDefault)
	query.AuthorFont = zen.Or(query.AuthorFont, authorFontDefault)
	query.WebsiteFont = zen.Or(query.WebsiteFont, websiteFontDefault)
	query.TitleFontSize = zen.Or(query.TitleFontSize, titleFontSizeDefault)
	query.AuthorFontSize = zen.Or(query.AuthorFontSize, authorFontSizeDefault)
	query.WebsiteFontSize = zen.Or(query.WebsiteFontSize, websiteFontSizeDefault)
	// Initialize image context
	img := gg.NewContext(query.Width, query.Height)
	// Background
	if query.Background != "" {
		// Load background
		bg, cleanup, err := extrender.LoadRemoteImage(query.Background)
		if err != nil {
			panic(err)
		}
		// Defer temp file cleanup
		defer cleanup()
		// Resize
		bg = imaging.Fill(bg, img.Width(), img.Height(), imaging.Center, imaging.Lanczos)
		// Write to image context
		img.DrawImage(bg, 0, 0)
	}
	// Overlay
	if query.BackgroundOverlay {
		// Define overlay position and size
		x := marginOverlay
		y := marginOverlay
		w := float64(img.Width()) - (2.0 * marginOverlay)
		h := float64(img.Height()) - (2.0 * marginOverlay)
		// Set overlay color
		img.SetColor(colorOverlay)
		// Draw overlay
		img.DrawRectangle(x, y, w, h)
		img.Fill()
	}
	// Dim
	if query.BackgroundDim != 0 {
		// Set dim color, depending on provided value
		img.SetColor(color.RGBA{0, 0, 0, uint8(query.BackgroundDim)})
		// Draw dim
		img.DrawRectangle(0, 0, float64(img.Width()), float64(img.Height()))
		img.Fill()
	}
	// Title
	if query.Title != "" {
		// Load title font
		if err := img.LoadFontFace(path.Join("dist/fonts", query.TitleFont), query.TitleFontSize); err != nil {
			panic("error while reading font file")
		}
		// Define title position and max width
		x := marginTitleX
		y := marginTitleY
		maxWidth := float64(img.Width()) - marginTitleX - marginTitleX
		// Define title color
		img.SetColor(colorTitle)
		// Draw title
		img.DrawStringWrapped(query.Title, x, y, 0, 0, maxWidth, 1.5, gg.AlignLeft)
	}
	// Author
	if query.Author != "" {
		// Load author font
		if err := img.LoadFontFace(path.Join("dist/fonts", query.AuthorFont), query.AuthorFontSize); err != nil {
			panic("error while reading font file")
		}
		// Define author position
		x := marginAuthorX
		y := float64(img.Height()) - marginAuthorY
		// Define author color
		img.SetColor(colorAuthor)
		// Draw author
		img.DrawString(query.Author, x, y)
	}
	// Website
	if query.Website != "" {
		// Load website font
		if err := img.LoadFontFace(path.Join("dist/fonts", query.WebsiteFont), query.WebsiteFontSize); err != nil {
			panic("error while reading font file")
		}
		// Define website position
		_, textHeight := img.MeasureString(query.Website)
		x := marginWebsiteX
		y := float64(img.Height()) - textHeight - marginWebsiteY
		// Define website color
		img.SetColor(colorWebsite)
		// Raw website
		img.DrawString(query.Website, x, y)
	}
	// Logo
	if query.Logo != "" {
		// Load logo
		logo, cleanup, err := extrender.LoadRemoteImage(query.Logo)
		if err != nil {
			panic(err)
		}
		// Defer temp file cleanup
		defer cleanup()
		// Resize
		logo = imaging.Resize(logo, 250, 0, imaging.Lanczos)
		// Define position
		x := float64(img.Width()) - float64(logo.Bounds().Dx()) - marginLogoX
		y := float64(img.Height()) - float64(logo.Bounds().Dy()) - marginLogoY
		// Write to image context
		img.DrawImage(logo, int(x), int(y))
	}
	// Generate unique og file
	ogfile := zen.Must(os.CreateTemp("/tmp", "*.oxigen.og"))
	// Defer clean up
	defer os.Remove(ogfile.Name())
	// Save resulting image to generated file
	zen.Must(0, img.SavePNG(ogfile.Name()))
	// Close file
	zen.Must(0, ogfile.Close())
	// Write response
	http.ServeFile(w, r, ogfile.Name())
}
