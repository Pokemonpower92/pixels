package imageset

import (
	"image"
	"image/color"
	"image/draw"
	"log"
)

func calculateAverageColors(images []*image.YCbCr) []color.RGBA {
	l := log.New(log.Writer(), "averageColor ", log.LstdFlags)
	l.Println("Calculating average colors.")

	var averageColors []color.RGBA
	for _, img := range images {
		b := img.Bounds()
		m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)

		bounds := m.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		// Initialize variables to accumulate color values
		var totalR, totalG, totalB, totalA uint32

		// Iterate through all pixels to calculate the sum of color values
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				pixelColor := m.At(x, y).(color.RGBA)
				totalR += uint32(pixelColor.R)
				totalG += uint32(pixelColor.G)
				totalB += uint32(pixelColor.B)
				totalA += uint32(pixelColor.A)
			}
		}

		// Calculate the average color values
		totalPixels := uint32(width * height)
		avgR := totalR / totalPixels
		avgG := totalG / totalPixels
		avgB := totalB / totalPixels
		avgA := totalA / totalPixels

		// Create the average color as a color.RGBA struct
		averageColor := color.RGBA{R: uint8(avgR), G: uint8(avgG), B: uint8(avgB), A: uint8(avgA)}

		averageColors = append(averageColors, averageColor)
	}

	l.Println("Average colors calculated.")

	return averageColors
}

type Generator struct {
	l     *log.Logger
	job   *Job
	store *S3Store
}

func NewGenerator(job *Job) *Generator {
	return &Generator{
		l:     log.New(log.Writer(), "generator ", log.LstdFlags),
		job:   job,
		store: NewS3Store(job.BucketName),
	}
}

func (g *Generator) Generate() (*ImageSet, error) {
	g.l.Printf("Generating imageset from job: %v", g.job)

	images, err := g.store.GetImageSet()
	if err != nil {
		g.l.Printf("Failed to get imageset from store: %s", err)
		return nil, err
	}

	im := &ImageSet{
		ID:            g.job.ImagesetID,
		Name:          g.job.BucketName,
		Description:   g.job.Description,
		AverageColors: calculateAverageColors(images),
	}
	return im, nil
}
