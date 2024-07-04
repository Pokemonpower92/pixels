package generator

import (
	"image"
	"image/color"
	"log"
	"strconv"

	"github.com/pokemonpower92/imagesetservice/internal/datastore"
	"github.com/pokemonpower92/imagesetservice/internal/domain"
	"github.com/pokemonpower92/imagesetservice/internal/job"
)

// calculateAverageColors calculates the average colors of a slice of images.
// It takes a slice of *image.RGBA as input and returns a slice of *color.RGBA.
func calculateAverageColors(images []*image.RGBA) []*color.RGBA {
	logger := log.New(log.Writer(), "averageColor ", log.LstdFlags)
	logger.Println("Calculating average colors.")

	var averageColors []*color.RGBA
	for _, image := range images {
		bounds := image.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		// Initialize variables to accumulate color values
		var totalRed, totalGreen, totalBlue, totalAlpha uint32

		// Iterate through all pixels to calculate the sum of color values
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				pixelColor := image.At(x, y).(color.RGBA)
				totalRed += uint32(pixelColor.R)
				totalGreen += uint32(pixelColor.G)
				totalBlue += uint32(pixelColor.B)
				totalAlpha += uint32(pixelColor.A)
			}
		}

		// Calculate the average color values
		totalPixels := uint32(width * height)
		avgRed := totalRed / totalPixels
		avgGreen := totalGreen / totalPixels
		avgBlue := totalBlue / totalPixels
		avgAlpha := totalAlpha / totalPixels

		// Create the average color as a color.RGBA struct
		averageColor := color.RGBA{R: uint8(avgRed), G: uint8(avgGreen), B: uint8(avgBlue), A: uint8(avgAlpha)}

		averageColors = append(averageColors, &averageColor)
	}

	logger.Println("Average colors calculated.")

	return averageColors
}

// Generator is an interface for generating image sets.
type Generator interface {
	Generate(job *job.Job) (*domain.ImageSet, error)
}

// ImageSetGenerator is a struct that implements the Generator interface.
type ImageSetGenerator struct {
	logger *log.Logger
	store  datastore.Store
}

// NewImageSetGenerator creates a new ImageSetGenerator instance.
// It takes a *Job as input and returns a pointer to ImageSetGenerator.
func NewImageSetGenerator(job *job.Job, logger *log.Logger) ImageSetGenerator {
	return ImageSetGenerator{
		logger: logger,
		store:  datastore.NewS3Store(job.BucketName),
	}
}

// Generate generates an image set based on the provided job.
// It takes a *Job as input and returns a pointer to types.ImageSet and an error.
func (generator ImageSetGenerator) Generate(job *job.Job) (*domain.ImageSet, error) {
	generator.logger.Printf("Generating imageset from job: %v", job)

	images, err := generator.store.GetImageSet()
	if err != nil {
		generator.logger.Printf("Failed to get imageset from store: %s", err)
		return nil, err
	}

	imagesetID, err := strconv.Atoi(job.ImagesetID)
	if err != nil {
		generator.logger.Printf("Failed to convert ImagesetID to int: %s", err)
		return nil, err
	}

	imageSet := &domain.ImageSet{
		ID:            imagesetID,
		Name:          job.BucketName,
		Description:   job.Description,
		AverageColors: calculateAverageColors(images),
	}
	return imageSet, nil
}
