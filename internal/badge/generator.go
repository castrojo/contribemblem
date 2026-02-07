package badge

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
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

	// Layout constants
	marginX        = 20
	marginTop      = 12
	statBarHeight  = 36
	accentHeight   = 3
	borderWidth    = 1
	statDividerW   = 1
	gradientStartX = 0.4 // Start gradient at 40% from left
)

var (
	// Destiny exotic gold
	PowerLevelColor = color.RGBA{244, 208, 63, 255} // #F4D03F
	AccentColor     = color.RGBA{244, 208, 63, 255} // #F4D03F
	WhiteColor      = color.RGBA{255, 255, 255, 255}
	DimWhiteColor   = color.RGBA{180, 180, 190, 255}
	BlackColor      = color.RGBA{0, 0, 0, 255}
	ShadowColor     = color.RGBA{0, 0, 0, 204} // alpha 0.8
	StatBarColor    = color.RGBA{0, 0, 0, 150}
	DividerColor    = color.RGBA{255, 255, 255, 50}
	BorderColor     = color.RGBA{60, 60, 65, 255}
	OverlayDark     = color.RGBA{0, 0, 0, 25}
)

// FontFaces holds the three font sizes used in badge generation
type FontFaces struct {
	Large  font.Face // 48pt for power level
	Medium font.Face // 26pt for username
	Small  font.Face // 14pt for stats
}

// Stats for badge generation
type Stats struct {
	Username     string
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

	// Phase 1: Scale and draw emblem as background
	xdraw.BiLinear.Scale(canvas, canvas.Bounds(), emblemImg, emblemImg.Bounds(), xdraw.Over, nil)

	// Phase 2: Overall darken overlay for Destiny dark UI feel
	drawRect(canvas, 0, 0, Width, Height, OverlayDark)

	// Phase 3: Horizontal gradient overlay (left transparent → right semi-opaque black)
	drawHorizontalGradient(canvas, int(float64(Width)*gradientStartX), 0, Width, Height, color.RGBA{0, 0, 0, 100})

	// Phase 4: Semi-transparent stat bar across bottom
	statBarY := Height - statBarHeight
	drawRect(canvas, 0, statBarY, Width, statBarHeight, StatBarColor)

	// Phase 5: Gold accent line at top
	drawRect(canvas, 0, 0, Width, accentHeight, AccentColor)

	// Phase 6: Border around entire badge
	drawBorder(canvas, Width, Height, borderWidth, BorderColor)

	// Load fonts
	fonts, err := loadFonts()
	if err != nil {
		return fmt.Errorf("failed to load fonts: %w", err)
	}
	defer fonts.Large.Close()
	defer fonts.Medium.Close()
	defer fonts.Small.Close()

	// Calculate Power Level
	powerLevel := stats.Commits + stats.PullRequests + stats.Issues + stats.Reviews + stats.Stars

	// Render username (top-left with medium font)
	if stats.Username != "" {
		usernameY := accentHeight + marginTop + 28
		DrawTextWithOutline(canvas, stats.Username, marginX+4, usernameY, fonts.Medium, WhiteColor)
	}

	// Render Power Level (right-aligned with diamond icon)
	powerText := fmt.Sprintf("%d", powerLevel)
	diamondIcon := "◆"

	// Measure text widths
	powerWidth := measureText(fonts.Large, powerText)
	diamondWidth := measureText(fonts.Large, diamondIcon)
	totalWidth := diamondWidth + 8 + powerWidth // 8px spacing between diamond and number

	// Calculate right-aligned position
	powerX := Width - marginX - totalWidth
	powerY := accentHeight + marginTop + 52

	// Draw diamond icon
	DrawTextWithOutline(canvas, diamondIcon, powerX, powerY, fonts.Large, PowerLevelColor)

	// Draw power level number
	DrawTextWithOutline(canvas, powerText, powerX+diamondWidth+8, powerY, fonts.Large, PowerLevelColor)

	// Render stats in stat bar (centered in cells with dividers)
	statLabels := []string{"COMMITS", "PRS", "ISSUES", "REVIEWS", "STARS"}
	statValues := []int{stats.Commits, stats.PullRequests, stats.Issues, stats.Reviews, stats.Stars}

	cellWidth := Width / 5
	statCenterY := statBarY + statBarHeight/2

	for i := 0; i < 5; i++ {
		// Draw vertical divider (except before first stat)
		if i > 0 {
			dividerX := i * cellWidth
			drawRect(canvas, dividerX, statBarY, statDividerW, statBarHeight, DividerColor)
		}

		// Prepare label and value text
		label := statLabels[i]
		value := FormatNumber(statValues[i])

		// Measure text widths for centering
		labelWidth := measureText(fonts.Small, label)
		valueWidth := measureText(fonts.Small, value)
		totalTextWidth := labelWidth + 6 + valueWidth // 6px spacing between label and value

		// Calculate center position within cell
		cellCenterX := i*cellWidth + cellWidth/2
		textStartX := cellCenterX - totalTextWidth/2

		// Draw label (dim white) and value (bright white)
		DrawTextWithOutline(canvas, label, textStartX, statCenterY+5, fonts.Small, DimWhiteColor)
		DrawTextWithOutline(canvas, value, textStartX+labelWidth+6, statCenterY+5, fonts.Small, WhiteColor)
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

func loadFonts() (*FontFaces, error) {
	ttf, err := opentype.Parse(rajdhaniFontData)
	if err != nil {
		return nil, err
	}

	large, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	medium, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    26,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		large.Close()
		return nil, err
	}

	small, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		large.Close()
		medium.Close()
		return nil, err
	}

	return &FontFaces{
		Large:  large,
		Medium: medium,
		Small:  small,
	}, nil
}

func savePNG(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// drawRect draws a filled rectangle with the given color
func drawRect(dst draw.Image, x, y, width, height int, col color.Color) {
	rect := image.Rect(x, y, x+width, y+height)
	draw.Draw(dst, rect, &image.Uniform{col}, image.Point{}, draw.Over)
}

// drawHorizontalGradient draws a left-to-right gradient from transparent to the given color
func drawHorizontalGradient(dst draw.Image, startX, y, endX, height int, endColor color.Color) {
	r, g, b, a := endColor.RGBA()
	maxAlpha := float64(a >> 8)

	for x := startX; x < endX; x++ {
		// Calculate alpha based on position (0.0 at startX, 1.0 at endX)
		progress := float64(x-startX) / float64(endX-startX)
		alpha := uint8(progress * maxAlpha)

		gradColor := color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: alpha,
		}

		drawRect(dst, x, y, 1, height, gradColor)
	}
}

// drawBorder draws a border around the image
func drawBorder(dst draw.Image, width, height, borderW int, col color.Color) {
	// Top
	drawRect(dst, 0, 0, width, borderW, col)
	// Bottom
	drawRect(dst, 0, height-borderW, width, borderW, col)
	// Left
	drawRect(dst, 0, 0, borderW, height, col)
	// Right
	drawRect(dst, width-borderW, 0, borderW, height, col)
}

// measureText returns the width of the text in pixels
func measureText(face font.Face, text string) int {
	drawer := &font.Drawer{
		Face: face,
	}

	advance := drawer.MeasureString(text)
	return advance.Round()
}
