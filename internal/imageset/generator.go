package imageset

import "log"

type Generator struct {
	l   *log.Logger
	job *Job
}

func NewGenerator(l *log.Logger, job *Job) *Generator {
	return &Generator{
		l:   l,
		job: job,
	}
}

func (g *Generator) Generate() (*ImageSet, error) {
	g.l.Printf("Generating imageset from job: %v", g.job)
	// Connect to image set providor
	// Load the image set.
	// Calculate the average color of the images
	// Return the image set
	im := &ImageSet{
		ID:          g.job.ImagesetID,
		Name:        "Test Image Set",
		Description: "This is a test image set",
		Images:      []string{"image1.jpg", "image2.jpg", "image3.jpg"},
		AverageColors: []float64{
			0.1,
			0.2,
			0.3,
		},
	}
	return im, nil
}
