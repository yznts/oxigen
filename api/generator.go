package api

import (
	"image/color"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/yuriizinets/oxigen/imgops"
	"go.kyoto.codes/zen/v3/errorsx"
	"go.kyoto.codes/zen/v3/httpx"
	"go.kyoto.codes/zen/v3/logic"
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
	marginWebsiteY = 50.0
	marginAuthorX  = 70.0
	marginAuthorY  = 150.0
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
	TitleFont       string  `query:"title.font"`
	TitleFontSize   float64 `query:"title.font.size"`
	AuthorFont      string  `query:"author.font"`
	AuthorFontSize  float64 `query:"author.font.size"`
	WebsiteFont     string  `query:"website.font"`
	WebsiteFontSize float64 `query:"website.font.size"`
	LogoScale       float64 `query:"logo.scale"`
	LogoAlignX      int     `query:"logo.align.x"`
	LogoAlignY      int     `query:"logo.align.y"`
	BackgroundDim   int     `query:"background.dim"`
	BackgroundFrame bool    `query:"background.overlay"`
}

func resolveHref(r *http.Request, href string) string {
	// If empty or schema provided, use as-is
	if href == "" || strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}
	// If only path provided, use referrer data to complete href
	if strings.HasPrefix(href, "/") {
		ref, _ := url.Parse(r.Referer()) // Parse referrer
		ref.Path = href                  // Override path
		return ref.String()
	}
	// If we're here, no schema provided (href starts with host).
	// So, we want to complete schema with request data.
	return r.Proto + "://" + href
}

func Generator(w http.ResponseWriter, r *http.Request) {
	// Unpack query
	query := GenerateQuery{}
	errorsx.Must(0, httpx.Query(r.URL.Query()).Unmarshal(&query))
	query.BackgroundFrame = r.URL.Query().Get("background.overlay") == "on"
	// Resolve defaults
	query.Width = logic.Or(query.Width, 1200)
	query.Height = logic.Or(query.Height, 628)
	query.Logo = resolveHref(r, query.Logo)
	query.Background = resolveHref(r, query.Background)
	query.TitleFont = logic.Or(query.TitleFont, titleFontDefault)
	query.AuthorFont = logic.Or(query.AuthorFont, authorFontDefault)
	query.WebsiteFont = logic.Or(query.WebsiteFont, websiteFontDefault)
	query.TitleFontSize = logic.Or(query.TitleFontSize, titleFontSizeDefault)
	query.AuthorFontSize = logic.Or(query.AuthorFontSize, authorFontSizeDefault)
	query.WebsiteFontSize = logic.Or(query.WebsiteFontSize, websiteFontSizeDefault)

	// Initialize image context
	img := gg.NewContext(query.Width, query.Height)
	// Background
	if query.Background != "" {
		// Load background
		bg, cleanup, err := imgops.GetRemote(query.Background)
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
	if query.BackgroundFrame {
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
		// Define title position and max width
		x := marginTitleX
		y := marginTitleY
		maxWidth := float64(img.Width()) - marginTitleX - marginTitleX
		// Draw title
		imgops.RenderText(img, imgops.Point{x, y}, imgops.Text{
			Text:  query.Title,
			Font:  query.TitleFont,
			Size:  query.TitleFontSize,
			Color: colorTitle,
			Width: maxWidth,
		})
	}
	// Author
	if query.Author != "" {
		// Define author position
		x := marginAuthorX
		y := float64(img.Height()) - marginOverlay - marginAuthorY
		// Draw author
		imgops.RenderText(img, imgops.Point{x, y}, imgops.Text{
			Text:  query.Author,
			Font:  query.AuthorFont,
			Size:  query.AuthorFontSize,
			Color: colorAuthor,
			Width: float64(img.Width()),
		})
	}
	// Website
	if query.Website != "" {
		// Define website position
		_, textHeight := img.MeasureString(query.Website)
		x := marginWebsiteX
		y := float64(img.Height()) - marginOverlay - textHeight - marginWebsiteY
		// Draw website
		imgops.RenderText(img, imgops.Point{x, y}, imgops.Text{
			Text:  query.Website,
			Font:  query.WebsiteFont,
			Size:  query.WebsiteFontSize,
			Color: colorWebsite,
			Width: float64(img.Width()),
		})
	}
	// Logo
	if query.Logo != "" {
		// Load logo
		logo, cleanup, err := imgops.GetRemote(query.Logo)
		if err != nil {
			panic(err)
		}
		// Defer temp file cleanup
		defer cleanup()
		// Resize
		if query.LogoScale != 0 {
			logo = imaging.Resize(logo, int(float64(logo.Bounds().Dx())*query.LogoScale), 0, imaging.Lanczos)
		}
		// Define position
		x := float64(img.Width()) - float64(logo.Bounds().Dx()) - marginLogoX
		y := float64(img.Height()) - float64(logo.Bounds().Dy()) - marginLogoY
		// Align
		x = x + float64(query.LogoAlignX)
		y = y + float64(query.LogoAlignY)
		// Write to image context
		img.DrawImage(logo, int(x), int(y))
	}
	// Generate unique og file
	ogfile := errorsx.Must(os.CreateTemp("/tmp", "*.oxigen.tmp"))
	// Defer clean up
	defer os.Remove(ogfile.Name())
	// Save resulting image to generated file
	errorsx.Must(0, img.SavePNG(ogfile.Name()))
	// Close file
	errorsx.Must(0, ogfile.Close())
	// Write response
	http.ServeFile(w, r, ogfile.Name())
}
