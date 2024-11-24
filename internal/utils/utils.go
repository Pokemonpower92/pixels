package utils

import (
	"image"
	"image/color"
	"image/draw"
	"io"
)

func ImageFileToRGBA(reader io.Reader) (*image.RGBA, error) {
	im, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	b := im.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), im, b.Min, draw.Src)
	return dst, nil
}

func CalculateAverageColor(image *image.RGBA) color.RGBA {
	r, g, b, a, totalPixels := 0, 0, 0, 0, 0
	bounds := image.Bounds()
	for y := bounds.Min.Y; y != bounds.Max.Y; y++ {
		for x := bounds.Min.X; x != bounds.Max.X; x++ {
			pixel := image.RGBAAt(x, y)
			r += int(pixel.R)
			g += int(pixel.G)
			b += int(pixel.B)
			a += int(pixel.A)
			totalPixels++
		}
	}
	return color.RGBA{
		R: uint8(r / totalPixels),
		G: uint8(g / totalPixels),
		B: uint8(b / totalPixels),
		A: uint8(a / totalPixels),
	}
}
