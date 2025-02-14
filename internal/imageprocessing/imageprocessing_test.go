package imageprocessing

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func createTestRGBAImage(fill color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	draw.Draw(
		img,
		img.Bounds(),
		&image.Uniform{fill},
		image.Point{},
		draw.Src,
	)
	return img
}

func TestCalculateAverageColor(t *testing.T) {
	testCases := []struct {
		name      string
		fillColor color.RGBA
		expected  color.RGBA
	}{
		{
			name:      "red",
			fillColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expected:  color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:      "green",
			fillColor: color.RGBA{R: 0, G: 255, B: 0, A: 255},
			expected:  color.RGBA{R: 0, G: 255, B: 0, A: 255},
		},
		{
			name:      "blue",
			fillColor: color.RGBA{R: 0, G: 0, B: 255, A: 255},
			expected:  color.RGBA{R: 0, G: 0, B: 255, A: 255},
		},
		{
			name:      "transparent",
			fillColor: color.RGBA{R: 0, G: 0, B: 0, A: 0},
			expected:  color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			img := createTestRGBAImage(test.fillColor)
			result := CalculateAverageColor(img)
			if result != test.expected {
				t.Errorf(
					"Test %s failed: expected %v got %v\n",
					test.name,
					test.expected,
					result,
				)
			}
		})
	}
}

func TestColorDifference(t *testing.T) {
	testCases := []struct {
		name        string
		firstColor  color.RGBA
		secondColor color.RGBA
		expected    float64
	}{
		{
			name:        "Same color",
			firstColor:  color.RGBA{R: 255, G: 0, B: 0, A: 255},
			secondColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expected:    0.0,
		},
	}
	for _, test := range testCases {
		result := CalculateColorDistance(test.firstColor, test.secondColor)
		if result != test.expected {
			t.Errorf(
				"Test %s failed: expected = %v, got %v",
				test.name,
				test.expected,
				result,
			)
		}
	}
}
