package generator

import (
	"image"
	"image/color"
	"log"
	"testing"

	"github.com/pokemonpower92/collagecommon/types"
)

func TestCalculateAverageColors(t *testing.T) {
	// Create test images
	img1 := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img1.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})     // Red
	img1.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})     // Green
	img1.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})     // Blue
	img1.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255}) // White

	img2 := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img2.Set(0, 0, color.RGBA{R: 128, G: 128, B: 128, A: 255}) // Gray
	img2.Set(1, 0, color.RGBA{R: 0, G: 0, B: 0, A: 255})       // Black
	img2.Set(0, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255}) // White
	img2.Set(1, 1, color.RGBA{R: 255, G: 0, B: 0, A: 255})     // Red

	images := []*image.RGBA{img1, img2}

	// Expected average colors
	expectedColors := []*color.RGBA{
		{R: 127, G: 127, B: 127, A: 255}, // Gray
		{R: 159, G: 95, B: 95, A: 255},   // Reddish
	}

	// Calculate average colors
	averageColors := calculateAverageColors(images)

	// Compare the calculated average colors with the expected colors
	for i, c := range averageColors {
		if c.R != expectedColors[i].R || c.G != expectedColors[i].G || c.B != expectedColors[i].B || c.A != expectedColors[i].A {
			t.Errorf("Average color mismatch at index %d. Expected: %v, Got: %v", i, expectedColors[i], c)
		}
	}
}

func TestGenerator_Generate(t *testing.T) {
	// Create a mock Job and Store
	job := &Job{
		ImagesetID:  "123",
		BucketName:  "test-bucket",
		Description: "Test description",
	}

	store := &MockStore{
		GetImageSetFunc: func() ([]*image.RGBA, error) {
			// Create a test image
			img := image.NewRGBA(image.Rect(0, 0, 2, 2))
			img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})     // Red
			img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})     // Green
			img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})     // Blue
			img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255}) // White

			return []*image.RGBA{img}, nil
		},
	}

	// Create a Generator instance
	g := &ImageSetGenerator{
		store:  store,
		logger: log.New(log.Writer(), "generator ", log.LstdFlags),
	}

	// Generate the ImageSet
	im, err := g.Generate(job)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the generated ImageSet
	expectedImageSet := &types.ImageSet{
		ID:            123,
		Name:          "test-bucket",
		Description:   "Test description",
		AverageColors: []*color.RGBA{{R: 127, G: 127, B: 127, A: 255}}, // Gray
	}

	if im.ID != expectedImageSet.ID || im.Name != expectedImageSet.Name || im.Description != expectedImageSet.Description {
		t.Errorf("ImageSet mismatch. Expected: %v, Got: %v", expectedImageSet, im)
	}

	if len(im.AverageColors) != len(expectedImageSet.AverageColors) {
		t.Errorf("Average colors length mismatch. Expected: %d, Got: %d", len(expectedImageSet.AverageColors), len(im.AverageColors))
	} else {
		for i, c := range im.AverageColors {
			if c.R != expectedImageSet.AverageColors[i].R || c.G != expectedImageSet.AverageColors[i].G || c.B != expectedImageSet.AverageColors[i].B || c.A != expectedImageSet.AverageColors[i].A {
				t.Errorf("Average color mismatch at index %d. Expected: %v, Got: %v", i, expectedImageSet.AverageColors[i], c)
			}
		}
	}
}

// MockStore is a mock implementation of the Store interface for testing purposes
type MockStore struct {
	GetImageSetFunc func() ([]*image.RGBA, error)
}

func (m *MockStore) GetImageSet() ([]*image.RGBA, error) {
	return m.GetImageSetFunc()
}
