package badge

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// DrawTextWithOutline renders multi-layer text (shadow → stroke → fill)
// Implements IMAGE-04 requirement for contrast on variable backgrounds
func DrawTextWithOutline(dst draw.Image, text string, x, y int, face font.Face, fillColor color.Color) {
	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}

	// Layer 1: Shadow (offset +2px x/y, black with alpha 0.8)
	shadowOffsets := []struct{ dx, dy int }{
		{2, 2}, {3, 3},
	}
	for _, offset := range shadowOffsets {
		shadowPoint := fixed.Point26_6{
			X: point.X + fixed.I(offset.dx),
			Y: point.Y + fixed.I(offset.dy),
		}
		drawer := &font.Drawer{
			Dst:  dst,
			Src:  image.NewUniform(ShadowColor),
			Face: face,
			Dot:  shadowPoint,
		}
		drawer.DrawString(text)
	}

	// Layer 2: Stroke (multi-offset technique, 4px effective width)
	strokeOffsets := []struct{ dx, dy int }{
		{-2, -2}, {-2, -1}, {-2, 0}, {-2, 1}, {-2, 2},
		{-1, -2}, {-1, 2},
		{0, -2}, {0, 2},
		{1, -2}, {1, 2},
		{2, -2}, {2, -1}, {2, 0}, {2, 1}, {2, 2},
	}
	for _, offset := range strokeOffsets {
		strokePoint := fixed.Point26_6{
			X: point.X + fixed.I(offset.dx),
			Y: point.Y + fixed.I(offset.dy),
		}
		drawer := &font.Drawer{
			Dst:  dst,
			Src:  image.NewUniform(BlackColor),
			Face: face,
			Dot:  strokePoint,
		}
		drawer.DrawString(text)
	}

	// Layer 3: Fill (main text color)
	drawer := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(fillColor),
		Face: face,
		Dot:  point,
	}
	drawer.DrawString(text)
}

// DrawTextWithTracking renders text with custom letter-spacing (tracking)
// Uses the same outline rendering as DrawTextWithOutline, but draws characters individually
func DrawTextWithTracking(dst draw.Image, text string, x, y int, face font.Face, fillColor color.Color, tracking int) {
	currentX := x
	for _, char := range text {
		charStr := string(char)
		point := fixed.Point26_6{
			X: fixed.I(currentX),
			Y: fixed.I(y),
		}

		// Layer 1: Shadow
		shadowOffsets := []struct{ dx, dy int }{
			{2, 2}, {3, 3},
		}
		for _, offset := range shadowOffsets {
			shadowPoint := fixed.Point26_6{
				X: point.X + fixed.I(offset.dx),
				Y: point.Y + fixed.I(offset.dy),
			}
			drawer := &font.Drawer{
				Dst:  dst,
				Src:  image.NewUniform(ShadowColor),
				Face: face,
				Dot:  shadowPoint,
			}
			drawer.DrawString(charStr)
		}

		// Layer 2: Stroke
		strokeOffsets := []struct{ dx, dy int }{
			{-2, -2}, {-2, -1}, {-2, 0}, {-2, 1}, {-2, 2},
			{-1, -2}, {-1, 2},
			{0, -2}, {0, 2},
			{1, -2}, {1, 2},
			{2, -2}, {2, -1}, {2, 0}, {2, 1}, {2, 2},
		}
		for _, offset := range strokeOffsets {
			strokePoint := fixed.Point26_6{
				X: point.X + fixed.I(offset.dx),
				Y: point.Y + fixed.I(offset.dy),
			}
			drawer := &font.Drawer{
				Dst:  dst,
				Src:  image.NewUniform(BlackColor),
				Face: face,
				Dot:  strokePoint,
			}
			drawer.DrawString(charStr)
		}

		// Layer 3: Fill
		drawer := &font.Drawer{
			Dst:  dst,
			Src:  image.NewUniform(fillColor),
			Face: face,
			Dot:  point,
		}
		drawer.DrawString(charStr)

		// Advance X by character width + tracking
		advance := drawer.MeasureString(charStr)
		currentX += advance.Round() + tracking
	}
}

// DrawTextSubtle renders text with lighter outline (1px instead of 2px)
// Use for username and secondary text that should feel integrated, not floating
func DrawTextSubtle(dst draw.Image, text string, x, y int, face font.Face, fillColor color.Color) {
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	// Layer 1: Shadow (lighter than full outline version)
	shadowPoint := fixed.Point26_6{X: point.X + fixed.I(1), Y: point.Y + fixed.I(2)}
	drawer := &font.Drawer{Dst: dst, Src: image.NewUniform(ShadowColor), Face: face, Dot: shadowPoint}
	drawer.DrawString(text)

	// Layer 2: Thin stroke (1px offsets instead of 2px)
	thinOffsets := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	for _, offset := range thinOffsets {
		strokePoint := fixed.Point26_6{X: point.X + fixed.I(offset.dx), Y: point.Y + fixed.I(offset.dy)}
		d := &font.Drawer{Dst: dst, Src: image.NewUniform(BlackColor), Face: face, Dot: strokePoint}
		d.DrawString(text)
	}

	// Layer 3: Fill
	d := &font.Drawer{Dst: dst, Src: image.NewUniform(fillColor), Face: face, Dot: point}
	d.DrawString(text)
}

// DrawTextWithGlow renders text with an outer glow effect followed by standard outline
// Use for power level numbers to create Destiny's luminous appearance
func DrawTextWithGlow(dst draw.Image, text string, x, y int, face font.Face, fillColor color.Color, glowColor color.Color) {
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	// Layer 0: Glow (low-alpha fill color at +-3px offsets)
	glowRGBA := glowColor.(color.RGBA)
	glowAlpha := color.RGBA{
		R: glowRGBA.R,
		G: glowRGBA.G,
		B: glowRGBA.B,
		A: 40, // very subtle
	}
	glowOffsets := []struct{ dx, dy int }{
		{-3, -3}, {-3, 0}, {-3, 3},
		{0, -3}, {0, 3},
		{3, -3}, {3, 0}, {3, 3},
		{-2, -2}, {-2, 2}, {2, -2}, {2, 2},
	}
	for _, offset := range glowOffsets {
		glowPoint := fixed.Point26_6{X: point.X + fixed.I(offset.dx), Y: point.Y + fixed.I(offset.dy)}
		d := &font.Drawer{Dst: dst, Src: image.NewUniform(glowAlpha), Face: face, Dot: glowPoint}
		d.DrawString(text)
	}

	// Then standard shadow + stroke + fill via DrawTextWithOutline
	DrawTextWithOutline(dst, text, x, y, face, fillColor)
}
