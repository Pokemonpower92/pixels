package imageset

import (
	"image/color"

	"github.com/pokemonpower92/collagecommon/types"
)

type DB interface {
	GetImageSet(id int) (*types.ImageSet, error)
	CreateImageSet(im *types.ImageSet) error
	SetAverageColors(id int, ac []*color.RGBA) error
}
