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
