package badge

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/fonts/Rajdhani-Bold.ttf
var rajdhaniFontData []byte

const (
	Width  = 800
	Height = 162 // Matches Destiny 2 emblem aspect ratio (474:96)
)

var (
	// Destiny exotic gold
	PowerLevelColor = color.RGBA{244, 208, 63, 255} // #F4D03F
	WhiteColor      = color.RGBA{255, 255, 255, 255}
	BlackColor      = color.RGBA{0, 0, 0, 255}
	ShadowColor     = color.RGBA{0, 0, 0, 204} // alpha 0.8
)

// Stats for badge generation
type Stats struct {
	Commits      int
	PullRequests int
	Issues       int
	Reviews      int
	Stars        int
}

// Generate creates badge image from emblem and stats
// emblemPath: path to emblem JPEG (data/emblem.jpg)
// stats: GitHub contribution stats
// outputPath: where to save badge PNG (badge.png)
func Generate(emblemPath string, stats *Stats, outputPath string) error {
	// Load emblem image
	emblemImg, err := loadImage(emblemPath)
	if err != nil {
		return fmt.Errorf("failed to load emblem: %w", err)
	}

	// Create canvas
	canvas := image.NewRGBA(image.Rect(0, 0, Width, Height))

	// Scale and draw emblem as background
	xdraw.BiLinear.Scale(canvas, canvas.Bounds(), emblemImg, emblemImg.Bounds(), xdraw.Over, nil)

	// Load font
	face48, face20, err := loadFonts()
	if err != nil {
		return fmt.Errorf("failed to load fonts: %w", err)
	}
	defer face48.Close()
	defer face20.Close()

	// Calculate Power Level
	powerLevel := stats.Commits + stats.PullRequests + stats.Issues + stats.Reviews + stats.Stars

	// Render Power Level (top-right)
	powerText := fmt.Sprintf("%d", powerLevel)
	DrawTextWithOutline(canvas, powerText, 700, 60, face48, PowerLevelColor)

	// Render individual stats (bottom row)
	statIcons := []string{"●", "◆", "■", "▲", "★"}
	statValues := []int{stats.Commits, stats.PullRequests, stats.Issues, stats.Reviews, stats.Stars}
	xOffset := 50
	spacing := 140

	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("%s %s", statIcons[i], FormatNumber(statValues[i]))
		DrawTextWithOutline(canvas, text, xOffset, 140, face20, WhiteColor)
		xOffset += spacing
	}

	// Save PNG
	return savePNG(canvas, outputPath)
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

func loadFonts() (font.Face, font.Face, error) {
	ttf, err := opentype.Parse(rajdhaniFontData)
	if err != nil {
		return nil, nil, err
	}

	face48, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, nil, err
	}

	face20, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, nil, err
	}

	return face48, face20, nil
}

func savePNG(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
