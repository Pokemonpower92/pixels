package imageset

import (
	"testing"

	"github.com/pokemonpower92/imagesetservice/config"
)

func TestNewS3Store(t *testing.T) {
	config.LoadEnvironmentVariables()

	// Create a new S3Store instance
	store := NewS3Store("mimikyu-collage-test")
	im, err := store.GetImageSet()
	t.Log(im, err)
}
