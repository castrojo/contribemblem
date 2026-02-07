package badge

import (
	"image"
	"image/color"
	"testing"

	"golang.org/x/image/font/opentype"
)

func TestDrawTextWithOutline(t *testing.T) {
	// Create an in-memory RGBA image
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))

	// Load the embedded Rajdhani font
	ttf, err := opentype.Parse(rajdhaniFontData)
	if err != nil {
		t.Fatalf("failed to parse embedded font: %v", err)
	}

	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: 0,
	})
	if err != nil {
		t.Fatalf("failed to create font face: %v", err)
	}
	defer face.Close()

	// Call DrawTextWithOutline with known text/position/color
	testText := "TEST"
	textX, textY := 50, 50
	fillColor := color.RGBA{255, 0, 0, 255} // Red

	DrawTextWithOutline(img, testText, textX, textY, face, fillColor)

	// Assert that pixels at the text position are modified (non-zero alpha)
	// Check a few pixels near the text position where we expect rendering
	textRegionPixels := []image.Point{
		{textX, textY - 5},      // Above baseline
		{textX + 10, textY - 5}, // Middle of first character
		{textX + 20, textY - 5}, // Second character
	}

	foundModified := false
	for _, p := range textRegionPixels {
		if p.X < 0 || p.Y < 0 || p.X >= img.Bounds().Dx() || p.Y >= img.Bounds().Dy() {
			continue
		}
		c := img.RGBAAt(p.X, p.Y)
		if c.A > 0 {
			foundModified = true
			break
		}
	}

	if !foundModified {
		t.Error("expected pixels near text position to be modified, but all had zero alpha")
	}

	// Assert pixels far from the text are unmodified (zero alpha)
	// Check corners and edges which should be untouched
	farPixels := []image.Point{
		{0, 0},                     // Top-left corner
		{img.Bounds().Dx() - 1, 0}, // Top-right corner
		{0, img.Bounds().Dy() - 1}, // Bottom-left corner
		{img.Bounds().Dx() - 1, img.Bounds().Dy() - 1}, // Bottom-right corner
	}

	for _, p := range farPixels {
		c := img.RGBAAt(p.X, p.Y)
		if c.A != 0 {
			t.Errorf("expected pixel at (%d,%d) to be unmodified (zero alpha), got alpha=%d", p.X, p.Y, c.A)
		}
	}

	// Additional validation: check that we have rendered multiple layers
	// Count distinct alpha values in the rendered region to verify shadow/stroke/fill layers
	alphaValues := make(map[uint8]bool)
	for y := textY - 20; y < textY+10; y++ {
		for x := textX - 10; x < textX+80; x++ {
			if x < 0 || y < 0 || x >= img.Bounds().Dx() || y >= img.Bounds().Dy() {
				continue
			}
			c := img.RGBAAt(x, y)
			if c.A > 0 {
				alphaValues[c.A] = true
			}
		}
	}

	// We expect at least 2 distinct alpha values (shadow layer has different alpha than fill)
	if len(alphaValues) < 2 {
		t.Errorf("expected at least 2 distinct alpha values from multi-layer rendering, got %d", len(alphaValues))
	}
}
