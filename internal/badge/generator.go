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
	"strings"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/fonts/Inter-Bold.ttf
var interBoldFontData []byte

//go:embed assets/fonts/Inter-Medium.ttf
var interMediumFontData []byte

const (
	Width  = 800
	Height = 162 // Matches Destiny 2 emblem aspect ratio (474:96)

	// Layout constants
	marginX        = 20
	marginTop      = 12
	statBarHeight  = 44 // Increased from 36 to fit value-over-label layout
	accentHeight   = 3
	borderWidth    = 1
	statDividerW   = 1
	gradientStartX = 0.4 // Start gradient at 40% from left
)

var (
	// Destiny exotic gold (warmer, more accurate to Destiny 2)
	PowerLevelColor = color.RGBA{245, 217, 106, 255} // #F5D96A - Destiny exotic gold
	AccentColor     = color.RGBA{206, 174, 51, 255}  // #CEAE33 - Destiny accent gold
	WhiteColor      = color.RGBA{255, 255, 255, 255}
	DimWhiteColor   = color.RGBA{180, 180, 190, 255}
	BlackColor      = color.RGBA{0, 0, 0, 255}
	ShadowColor     = color.RGBA{0, 0, 0, 204}      // alpha 0.8
	StatBarColor    = color.RGBA{0, 0, 0, 170}      // Increased from 150 for better readability
	DividerColor    = color.RGBA{255, 255, 255, 70} // Increased from 50 for better visibility
	BorderColor     = color.RGBA{45, 45, 50, 255}
	OverlayDark     = color.RGBA{0, 0, 0, 35}
)

// FontFaces holds the four font faces used in badge generation
type FontFaces struct {
	Large     font.Face // 48pt Inter Bold - power level
	Medium    font.Face // 20pt Inter Medium - username
	StatValue font.Face // 16pt Inter Bold - stat numbers
	StatLabel font.Face // 10pt Inter Medium - stat labels (ALL-CAPS)
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

	// Phase 1: Scale and draw emblem as background with aspect-ratio-aware crop-to-fill
	// Calculate target aspect ratio (800:162 = 4.94:1)
	targetAspect := float64(Width) / float64(Height)
	srcBounds := emblemImg.Bounds()
	srcWidth := float64(srcBounds.Dx())
	srcHeight := float64(srcBounds.Dy())
	srcAspect := srcWidth / srcHeight

	// Determine crop rectangle that matches target aspect ratio
	var cropRect image.Rectangle
	if srcAspect > targetAspect {
		// Source is wider - crop horizontally (center crop)
		newWidth := int(srcHeight * targetAspect)
		offsetX := (srcBounds.Dx() - newWidth) / 2
		cropRect = image.Rect(
			srcBounds.Min.X+offsetX,
			srcBounds.Min.Y,
			srcBounds.Min.X+offsetX+newWidth,
			srcBounds.Max.Y,
		)
	} else {
		// Source is taller - crop vertically (center crop)
		newHeight := int(srcWidth / targetAspect)
		offsetY := (srcBounds.Dy() - newHeight) / 2
		cropRect = image.Rect(
			srcBounds.Min.X,
			srcBounds.Min.Y+offsetY,
			srcBounds.Max.X,
			srcBounds.Min.Y+offsetY+newHeight,
		)
	}

	// Scale the cropped region to fill the canvas
	xdraw.BiLinear.Scale(canvas, canvas.Bounds(), emblemImg, cropRect, xdraw.Over, nil)

	// Phase 2: Overall darken overlay for Destiny dark UI feel
	drawRect(canvas, 0, 0, Width, Height, OverlayDark)

	// Phase 3: Horizontal gradient overlay (left transparent → right semi-opaque black)
	drawHorizontalGradient(canvas, int(float64(Width)*gradientStartX), 0, Width, Height, color.RGBA{0, 0, 0, 150})

	// Phase 3b: Bottom vignette (subtle bottom-up darkening above stat bar)
	vignetteStartY := Height / 2           // start at vertical midpoint
	vignetteEndY := Height - statBarHeight // end at stat bar top
	drawVerticalGradient(canvas, 0, vignetteStartY, Width, vignetteEndY, color.RGBA{0, 0, 0, 60})

	// Phase 4: Semi-transparent stat bar across bottom
	statBarY := Height - statBarHeight
	drawRect(canvas, 0, statBarY, Width, statBarHeight, StatBarColor)

	// Stat bar top edge separator
	StatBarEdgeColor := color.RGBA{255, 255, 255, 30} // very subtle white line
	drawRect(canvas, 0, statBarY, Width, 1, StatBarEdgeColor)

	// Phase 5: Gold accent line at top
	drawRect(canvas, 0, 0, Width, accentHeight, AccentColor)
	// Accent glow (subtle bloom below the solid line)
	AccentGlowColor := color.RGBA{206, 174, 51, 80} // translucent gold
	drawRect(canvas, 0, accentHeight, Width, 1, AccentGlowColor)

	// Phase 6: Border around entire badge
	drawBorder(canvas, Width, Height, borderWidth, BorderColor)

	// Load fonts
	fonts, err := loadFonts()
	if err != nil {
		return fmt.Errorf("failed to load fonts: %w", err)
	}
	defer fonts.Large.Close()
	defer fonts.Medium.Close()
	defer fonts.StatValue.Close()
	defer fonts.StatLabel.Close()

	// Calculate Power Level
	powerLevel := stats.Commits + stats.PullRequests + stats.Issues + stats.Reviews + stats.Stars

	// Render username (positioned to the right of emblem icon, centered in left portion)
	// In Destiny 2, usernames appear around 370-400px from left edge on 800px canvas
	if stats.Username != "" {
		usernameY := accentHeight + marginTop + 22 // Adjusted for smaller 20pt font
		usernameX := 130                           // Positioned right of emblem icon area (~80-120px)
		DrawTextSubtle(canvas, strings.ToUpper(stats.Username), usernameX, usernameY, fonts.Medium, WhiteColor)
	}

	// Render Power Level (right-aligned with programmatic diamond icon)
	powerText := fmt.Sprintf("%d", powerLevel)

	// Diamond sizing: proper diamond proportions (equal width and height)
	// Increased slightly for better visibility next to 48pt text
	diamondHalfH := 10 // 20px total height
	diamondHalfW := 10 // 20px total width (equal to height for proper diamond)
	diamondGap := 6    // gap between diamond right edge and number left edge

	// Power level number position (right-aligned)
	powerWidth := measureText(fonts.Large, powerText)
	totalWidth := (diamondHalfW * 2) + diamondGap + powerWidth
	powerX := Width - marginX - totalWidth
	powerY := accentHeight + marginTop + 52

	// Diamond center position (vertically centered with number baseline)
	// The baseline is at powerY, cap height extends upward ~35px
	// Center the diamond vertically with the number: baseline - capHeight/2
	diamondCX := powerX + diamondHalfW
	diamondCY := powerY - 16 // adjust to visually center with number

	// Draw diamond with outline for contrast (same 3-layer approach as text)
	// Layer 1: shadow
	drawDiamond(canvas, diamondCX+2, diamondCY+2, diamondHalfW+1, diamondHalfH+1, ShadowColor)
	// Layer 2: black outline
	drawDiamond(canvas, diamondCX, diamondCY, diamondHalfW+2, diamondHalfH+2, BlackColor)
	// Layer 3: gold fill
	drawDiamond(canvas, diamondCX, diamondCY, diamondHalfW, diamondHalfH, PowerLevelColor)

	// Draw power level number after diamond with glow effect
	DrawTextWithGlow(canvas, powerText, powerX+(diamondHalfW*2)+diamondGap, powerY, fonts.Large, PowerLevelColor, PowerLevelColor)

	// Render stats in stat bar (vertical value-over-label layout)
	statLabels := []string{"COMMITS", "PRS", "ISSUES", "REVIEWS", "STARS"}
	statValues := []int{stats.Commits, stats.PullRequests, stats.Issues, stats.Reviews, stats.Stars}

	cellWidth := Width / 5

	for i := 0; i < 5; i++ {
		// Draw vertical divider (except before first stat)
		if i > 0 {
			dividerX := i * cellWidth
			drawRect(canvas, dividerX, statBarY, statDividerW, statBarHeight, DividerColor)
		}

		label := statLabels[i]
		value := FormatNumber(statValues[i])

		// Center value horizontally in cell
		valueWidth := measureText(fonts.StatValue, value)
		cellCenterX := i*cellWidth + cellWidth/2
		valueX := cellCenterX - valueWidth/2

		// Center label horizontally in cell
		labelWidth := measureText(fonts.StatLabel, label)
		labelX := cellCenterX - labelWidth/2

		// Value on upper line: 18px from stat bar top
		valueY := statBarY + 18
		DrawTextWithOutline(canvas, value, valueX, valueY, fonts.StatValue, WhiteColor)

		// Label on lower line: 36px from stat bar top
		labelY := statBarY + 36
		DrawTextWithOutline(canvas, label, labelX, labelY, fonts.StatLabel, DimWhiteColor)
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
	// Parse Inter Bold font
	boldTTF, err := opentype.Parse(interBoldFontData)
	if err != nil {
		return nil, err
	}

	// Parse Inter Medium font
	mediumTTF, err := opentype.Parse(interMediumFontData)
	if err != nil {
		return nil, err
	}

	// Large: 48pt Inter Bold - power level
	large, err := opentype.NewFace(boldTTF, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	// Medium: 20pt Inter Medium - username
	medium, err := opentype.NewFace(mediumTTF, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		large.Close()
		return nil, err
	}

	// StatValue: 16pt Inter Bold - stat numbers
	statValue, err := opentype.NewFace(boldTTF, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		large.Close()
		medium.Close()
		return nil, err
	}

	// StatLabel: 10pt Inter Medium - stat labels (ALL-CAPS)
	statLabel, err := opentype.NewFace(mediumTTF, &opentype.FaceOptions{
		Size:    10,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		large.Close()
		medium.Close()
		statValue.Close()
		return nil, err
	}

	return &FontFaces{
		Large:     large,
		Medium:    medium,
		StatValue: statValue,
		StatLabel: statLabel,
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

// drawVerticalGradient draws a top-to-bottom gradient from transparent to the given color
func drawVerticalGradient(dst draw.Image, x, startY, width, endY int, endColor color.Color) {
	r, g, b, a := endColor.RGBA()
	maxAlpha := float64(a >> 8)

	for y := startY; y < endY; y++ {
		progress := float64(y-startY) / float64(endY-startY)
		alpha := uint8(progress * maxAlpha)
		gradColor := color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: alpha,
		}
		drawRect(dst, x, y, width, 1, gradColor)
	}
}

// drawDiamond draws a filled diamond (rotated square) centered at (cx, cy)
// with the given half-width and half-height, filled with the specified color.
func drawDiamond(dst draw.Image, cx, cy, halfW, halfH int, col color.Color) {
	// Fill a diamond by iterating rows from top to bottom
	// At each row y, calculate the horizontal span using linear interpolation
	for dy := -halfH; dy <= halfH; dy++ {
		// Calculate width at this row (full width at center, 0 at tips)
		progress := 1.0 - float64(abs(dy))/float64(halfH)
		spanHalf := int(float64(halfW) * progress)
		for dx := -spanHalf; dx <= spanHalf; dx++ {
			dst.Set(cx+dx, cy+dy, col)
		}
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
