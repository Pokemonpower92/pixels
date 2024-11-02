package service

import (
	"fmt"

	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

func CreateCollage(collage *sqlc.Collage) {
	fmt.Printf("Creating collage: %+v\n", collage)
	// Find all of the images and where they will go
	// Do this concurrently for each section of the
	// target image.
	// Store the results in the collage_images table.

	// Construct the collage.
}

func findAverageColorForSection() {
	fmt.Printf("Finding average color for section.")
}

func determineImagePlacements() {
	fmt.Printf("Finding average color for section.")
}
