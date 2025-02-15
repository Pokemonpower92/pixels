package service

import (
	"bytes"
	"encoding/json"
	"image"
	"image/draw"
	"image/png"
	"sync"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

// GenerateCollage generates the final image for a collage.
func GenerateCollage(collageImage *sqlc.CollageImage) {
	logger := NewServiceLogger("CollageGenerator")
	store := datastore.NewStore()
	generator := CollageGenerator{
		logger:         logger,
		store:          store,
		collageImageId: collageImage.ID,
		collageId:      collageImage.CollageID,
	}
	generator.generate()
}

type CollageGenerator struct {
	logger         *ServiceLogger
	store          datastore.Store
	collageImageId uuid.UUID
	collageId      uuid.UUID
}

func (cg *CollageGenerator) generate() {
	cg.logger.Printf("Generating collage for collage %s\n", cg.collageId)
	// Get metadata file as collageMetaData
	metaData, err := cg.getMetaData()
	if err != nil {
		cg.logger.Printf("Error getting collage metadata: %v\n", err)
		return
	}
	// Create blank canvas.
	canvas, err := cg.createCanvas(metaData.Resolution)
	if err != nil {
		cg.logger.Printf("Error getting creating canvas: %v\n", err)
		return
	}
	// Build image concurrently.
	if err := cg.buildFinalImage(canvas, metaData); err != nil {
		cg.logger.Printf("Error generating final image: %s\n", err)
	}
	cg.logger.Printf("Final image generated\n")
}

func (cg *CollageGenerator) getMetaData() (*CollageMetaData, error) {
	metaDataFile, err := cg.store.GetFile(cg.collageId)
	if err != nil {
		return nil, err
	}
	var metaData CollageMetaData
	decoder := json.NewDecoder(metaDataFile)
	decoder.Decode(&metaData)
	return &metaData, nil
}

func (cg *CollageGenerator) createCanvas(
	resolution config.ResolutionConfig,
) (*image.RGBA, error) {
	canvas := image.NewRGBA(image.Rect(
		0,
		0,
		resolution.CollageWidth,
		resolution.CollageHeight,
	))
	return canvas, nil
}

func (cg *CollageGenerator) fillSections(
	start int,
	chunkSize int,
	canvas *image.RGBA,
	metaData *CollageMetaData,
) {
	for section := start; section < start+chunkSize; section++ {
		// Load image for section.
		fileId := metaData.SectionMap[section]
		imageFile, err := cg.store.GetFile(fileId)
		if err != nil {
			cg.logger.Printf(
				"Error getting fill image: %v\n",
				err,
			)
		}
		im, _, err := image.Decode(imageFile)
		if err != nil {
			cg.logger.Printf(
				"Error decoding fill image: %v\n",
				err,
			)
		}
		// Scale it.
		sectionWidth := metaData.Resolution.SectionWidth
		sectionHeight := metaData.Resolution.SectionHeight
		scaledImage := resize.Resize(
			uint(sectionWidth),
			uint(sectionHeight),
			im,
			resize.Lanczos2,
		)
		// Draw it where it needs to go on the canvas.
		numColumns := metaData.Resolution.CollageWidth / sectionWidth
		row := section / numColumns
		col := section % numColumns
		startingPoint := image.Point{
			X: sectionWidth * col,
			Y: sectionHeight * row,
		}
		bounds := image.Rectangle{
			Min: startingPoint,
			Max: startingPoint.Add(scaledImage.Bounds().Size()),
		}
		draw.Draw(
			canvas,
			bounds,
			scaledImage,
			scaledImage.Bounds().Min,
			draw.Src,
		)
	}
}

func (cg *CollageGenerator) buildFinalImage(
	canvas *image.RGBA,
	metaData *CollageMetaData,
) error {
	numThreads := 10
	chunkSize := len(metaData.SectionMap) / numThreads
	var wg sync.WaitGroup
	for thread := 0; thread < numThreads; thread++ {
		wg.Add(1)
		startingSection := thread * chunkSize
		go func(start int) {
			defer wg.Done()
			cg.fillSections(
				start,
				chunkSize,
				canvas,
				metaData,
			)
		}(startingSection)
	}
	wg.Wait()
	var buf bytes.Buffer
	err := png.Encode(&buf, canvas)
	if err != nil {
		return err
	}
	cg.store.PutFile(cg.collageImageId, &buf)
	return nil
}
