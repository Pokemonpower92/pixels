package generator

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"testing"

	"github.com/pokemonpower92/collagegenerator/internal/domain"
	"github.com/pokemonpower92/collagegenerator/internal/job"
)

type mockDatastore struct{}

func (m *mockDatastore) GetImages() ([]*image.RGBA, error) {
	red := color.RGBA{R: 255, G: 0, B: 0, A: 0}
	green := color.RGBA{R: 0, G: 255, B: 0, A: 0}
	blue := color.RGBA{R: 0, G: 0, B: 255, A: 0}

	redImage := image.NewRGBA(image.Rect(0, 0, 2, 2))
	greenImage := image.NewRGBA(image.Rect(0, 0, 2, 2))
	blueImage := image.NewRGBA(image.Rect(0, 0, 2, 2))

	draw.Draw(redImage, redImage.Bounds(), &image.Uniform{red}, image.Point{}, draw.Src)
	draw.Draw(greenImage, greenImage.Bounds(), &image.Uniform{green}, image.Point{}, draw.Src)
	draw.Draw(blueImage, blueImage.Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)

	return []*image.RGBA{redImage, greenImage, blueImage}, nil
}

func TestCalculateAverageColors(t *testing.T) {
	tests := []struct {
		name           string
		images         []*image.RGBA
		expectedColors []*color.RGBA
	}{
		{
			name: "Test case 1",
			images: []*image.RGBA{
				image.NewRGBA(image.Rect(0, 0, 2, 2)),
				image.NewRGBA(image.Rect(0, 0, 2, 2)),
			},
			expectedColors: []*color.RGBA{
				{R: 0, G: 0, B: 0, A: 0},
				{R: 0, G: 0, B: 0, A: 0},
			},
		},
		// Add more test cases here...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			averageColors := calculateAverageColors(tt.images)

			if len(averageColors) != len(tt.expectedColors) {
				t.Errorf("Average colors length mismatch. Expected: %d, Got: %d",
					len(tt.expectedColors),
					len(averageColors),
				)
			} else {
				for i, c := range averageColors {
					if c.R != tt.expectedColors[i].R ||
						c.G != tt.expectedColors[i].G ||
						c.B != tt.expectedColors[i].B ||
						c.A != tt.expectedColors[i].A {
						t.Errorf(
							"Average color mismatch at index %d. Expected: %v, Got: %v",
							i,
							tt.expectedColors[i],
							c,
						)
					}
				}
			}
		})
	}
}

func TestGenerator(t *testing.T) {
	tests := []struct {
		name     string
		input    *job.ImageSetJob
		expected *domain.ImageSet
	}{
		{
			name: "Test case 1",
			input: &job.ImageSetJob{
				ImagesetID:  "1",
				BucketName:  "Test case 1",
				Description: "Test case 1",
			},
			expected: &domain.ImageSet{
				ID:          1,
				Name:        "Test case 1",
				Description: "Test case 1",
				AverageColors: []*color.RGBA{
					// Red
					{R: 255, G: 0, B: 0, A: 0},
					// Green
					{R: 0, G: 255, B: 0, A: 0},
					// Blue
					{R: 0, G: 0, B: 255, A: 0},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isg := ImageSetGenerator{
				logger: log.New(log.Writer(), "testLogger: ", log.LstdFlags),
				store:  &mockDatastore{},
			}
			result, err := isg.Generate(tt.input)
			if err != nil {
				t.Errorf("Error generating image set: %s", err)
			}
			if result.ID != tt.expected.ID {
				t.Errorf(
					"ImageSet ID mismatch. Expected: %d, Got: %d",
					tt.expected.ID,
					result.ID,
				)
			}
			if result.Name != tt.expected.Name {
				t.Errorf(
					"ImageSet Name mismatch. Expected: %s, Got: %s",
					tt.expected.Name,
					result.Name,
				)
			}
			if result.AverageColors[0].R != tt.expected.AverageColors[0].R ||
				result.AverageColors[0].G != tt.expected.AverageColors[0].G ||
				result.AverageColors[0].B != tt.expected.AverageColors[0].B ||
				result.AverageColors[0].A != tt.expected.AverageColors[0].A {
				t.Errorf("Average color mismatch at index 0. Expected: %v, Got: %v",
					tt.expected.AverageColors[0],
					result.AverageColors[0],
				)
			}
		})
	}
}
