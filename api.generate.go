package main

import (
	"image/color"
	"io/ioutil"
	"net/http"
	"os"

	"git.sr.ht/~kyoto-framework/zen"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

// Generation constants
var (
	fontTitle   = "assets/fonts/OpenSans/OpenSans-SemiBold.ttf"
	fontAuthor  = "assets/fonts/OpenSans/OpenSans-SemiBold.ttf"
	fontWebsite = "assets/fonts/OpenSans/OpenSans-Light.ttf"

	fontSizeTitle   = 80.0
	fontSizeAuthor  = 50.0
	fontSizeWebsite = 50.0

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
	// Generation parameters
	Width  int `query:"width"`
	Height int `query:"height"`
	// Data
	Title   string `query:"title"`
	Author  string `query:"author"`
	Website string `query:"website"`
	Logo    string `query:"logo"`
	// Settings
	Background string `query:"background"`
	Overlay    bool   `query:"overlay"`
	Dim        int    `query:"dim"`
}

func AGenerate(w http.ResponseWriter, r *http.Request) {
	// Unpack query
	query := GenerateQuery{}
	zen.Must(0, zen.Query(r.URL.Query()).Unmarshal(&query))
	// Initialize image context
	img := gg.NewContext(
		zen.Or(query.Width, 1200),
		zen.Or(query.Height, 628),
	)
	// Background
	if query.Background != "" {
		// Download background from provided url
		bgresp := zen.Request("GET", query.Background).Do()
		bgbytes := zen.Must(ioutil.ReadAll(bgresp.Body))
		bgresp.Body.Close()
		// Open temp file
		bgfile := zen.Must(ioutil.TempFile("/tmp", "*.oxigen.bg"))
		defer os.Remove(bgfile.Name())
		// Save loaded image to temp file
		zen.Must(bgfile.Write(bgbytes))
		// Close file
		zen.Must(0, bgfile.Close())
		// Load into object
		bg := zen.Must(gg.LoadImage(bgfile.Name()))
		// Resize
		bg = imaging.Fill(bg, img.Width(), img.Height(), imaging.Center, imaging.Lanczos)
		// Write to image context
		img.DrawImage(bg, 0, 0)
	}
	// Overlay
	if query.Overlay {
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
	if query.Dim != 0 {
		// Set dim color, depending on provided value
		img.SetColor(color.RGBA{0, 0, 0, uint8(query.Dim)})
		// Draw dim
		img.DrawRectangle(0, 0, float64(img.Width()), float64(img.Height()))
		img.Fill()
	}
	// Title
	if query.Title != "" {
		// Load title font
		if err := img.LoadFontFace(fontTitle, fontSizeTitle); err != nil {
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
		if err := img.LoadFontFace(fontAuthor, fontSizeAuthor); err != nil {
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
		if err := img.LoadFontFace(fontWebsite, fontSizeWebsite); err != nil {
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
		// Download logo from provided url
		lgresp := zen.Request("GET", query.Logo).Do()
		lgbytes := zen.Must(ioutil.ReadAll(lgresp.Body))
		lgresp.Body.Close()
		// Open temp file
		lgfile, _ := ioutil.TempFile("/tmp", "*.oxigen.lg")
		defer lgfile.Close()
		// Save loaded image to temp file
		zen.Must(lgfile.Write(lgbytes))
		// Close file
		zen.Must(0, lgfile.Close())
		// Load into object
		bg := zen.Must(gg.LoadImage(lgfile.Name()))
		// Resize
		bg = imaging.Resize(bg, 250, 0, imaging.Lanczos)
		// Define position
		x := float64(img.Width()) - float64(bg.Bounds().Dx()) - marginLogoX
		y := float64(img.Height()) - float64(bg.Bounds().Dy()) - marginLogoY
		// Write to image context
		img.DrawImage(bg, int(x), int(y))
	}
	// Generate unique og file name
	ogfile := zen.Must(ioutil.TempFile("/tmp", "*.oxigen.og")).Name()
	// Save resulting image to generated og name
	zen.Must(0, img.SavePNG(ogfile))
	// Write response
	http.ServeFile(w, r, ogfile)
}
