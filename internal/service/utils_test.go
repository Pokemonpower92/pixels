package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/pokemonpower92/collagegenerator/config"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
	"github.com/pokemonpower92/collagegenerator/internal/stubs"
)

type ACRepoExtenderStub struct {
	GetByImageSetFunc func(id uuid.UUID) ([]*sqlc.AverageColor, error)
}

func (acr *ACRepoExtenderStub) GetByImageSetId(id uuid.UUID) ([]*sqlc.AverageColor, error) {
	return acr.GetByImageSetFunc(id)
}

func successRepo() ACRepoExtenderStub {
	return ACRepoExtenderStub{
		GetByImageSetFunc: func(id uuid.UUID) ([]*sqlc.AverageColor, error) {
			return []*sqlc.AverageColor{
				{
					DbID:     1,
					ID:       uuid.New(),
					FileName: "stubFile",
					R:        1,
					G:        1,
					B:        1,
					A:        1,
				},
			}, nil
		},
	}
}

func errorRepo() ACRepoExtenderStub {
	return ACRepoExtenderStub{
		GetByImageSetFunc: func(id uuid.UUID) ([]*sqlc.AverageColor, error) {
			return nil, errors.New("Stub error.")
		},
	}
}

func metaDataReader() io.Reader {
	metaData := CollageMetaData{
		Resolution: config.ResolutionConfig{},
		SectionMap: make([]uuid.UUID, 0),
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.Encode(metaData)
	return &buf
}

func imageReader() io.Reader {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.RGBA{R: 100, G: 100, B: 100, A: 255})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	imageBytes := buf.Bytes()
	return bytes.NewReader(imageBytes)
}

func textReader() io.Reader {
	return strings.NewReader("Hello, Reader!")
}

func successStore() stubs.StoreStub {
	return stubs.StoreStub{
		GetRGBAFunc: func(id uuid.UUID) (*image.RGBA, error) {
			return nil, nil
		},
		GetFileFunc: func(id uuid.UUID) (io.Reader, error) {
			// Return a new reader each time to avoid EOF issues
			return imageReader(), nil
		},
		PutFileFunc: func(id uuid.UUID, reader io.Reader) error {
			return nil
		},
	}
}
