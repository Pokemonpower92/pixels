package imageprocessing

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
)

// 16-bit RGBA for more fidelity.
type RGB16 struct {
	R, G, B uint16
}

func (c RGB16) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	g = uint32(c.G)
	b = uint32(c.B)
	a = 0xffff
	return
}

// Calculate the average color of the given image.RGBA
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

// Calculate the difference between color.RGBAs c1 and c2
// Difference is weighted by human perception.
func CalculateColorDistance(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	rf1 := float64(r1)
	gf1 := float64(g1)
	bf1 := float64(b1)

	rf2 := float64(r2)
	gf2 := float64(g2)
	bf2 := float64(b2)

	r := rf1 - rf2
	g := gf1 - gf2
	b := bf1 - bf2

	return math.Sqrt(r*r + g*g + b*b)
}
