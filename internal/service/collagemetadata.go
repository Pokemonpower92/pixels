package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/filestore"
	"github.com/pokemonpower92/collagegenerator/internal/imageprocessing"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type CollageMetaData struct {
	Resolution config.ResolutionConfig `json:"resolution"`
	SectionMap []uuid.UUID             `json:"section_map"`
}

func CreateCollageMetaData(collage *sqlc.Collage, logger *slog.Logger) {
	serviceContext, cancel := context.WithTimeout(
		context.Background(),
		time.Second*30,
	)
	defer cancel()
	postgresConfig := config.NewPostgresConfig()
	acRepo, err := repository.NewAverageColorRepository(
		postgresConfig,
		serviceContext,
	)
	if err != nil {
		panic("Couldn't create repo")
	}
	store := filestore.NewStore(logger)
	service := newCollageMetaDataService(
		collage,
		acRepo,
		store,
		logger,
	)
	service.determineImagePlacements()
}

type collageMetaDataService struct {
	logger     *slog.Logger
	numThreads int
	collage    *sqlc.Collage
	acRepo     repository.ACRepo
	resolution *config.ResolutionConfig
	sectionMap []uuid.UUID
	store      filestore.Store
}

func newCollageMetaDataService(
	collage *sqlc.Collage,
	acRepo repository.ACRepo,
	store filestore.Store,
	logger *slog.Logger,
) *collageMetaDataService {
	resolution := config.NewResolutionConfig()
	sectionMap := make(
		[]uuid.UUID,
		resolution.XSections*resolution.YSections,
	)
	return &collageMetaDataService{
		numThreads: 10,
		collage:    collage,
		acRepo:     acRepo,
		resolution: resolution,
		sectionMap: sectionMap,
		store:      store,
		logger:     logger,
	}
}

func (cs *collageMetaDataService) getAverageColors() ([]*sqlc.AverageColor, error) {
	averageColors, err := cs.acRepo.GetByResourceId(cs.collage.ImageSetID)
	if err != nil {
		return nil, errors.New("Failed to get average colors")
	}
	return averageColors, nil
}

// Get the local average color value of the collage's
// target image by scaling it down to X_SECTIONSxY_SECTIONS
func (cs *collageMetaDataService) getSectionAverageColors() ([]color.Color, error) {
	targetImageReader, err := cs.store.GetFile(cs.collage.TargetImageID)
	if err != nil {
		cs.logger.Error(fmt.Sprintf(
			"Failed to load target image: %s\n",
			cs.collage.TargetImageID,
		))
		return nil, err
	}
	targetImage, _, err := image.Decode(targetImageReader)
	if err != nil {
		cs.logger.Error(fmt.Sprintf(
			"Failed to decode target image: %s\n",
			cs.collage.TargetImageID,
		))
		return nil, err
	}
	scaledImage := resize.Resize(
		uint(cs.resolution.XSections),
		uint(cs.resolution.YSections),
		targetImage,
		resize.Lanczos2,
	)
	bounds := scaledImage.Bounds()
	averageColors := make([]color.Color, bounds.Dx()*bounds.Dy())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			averageColors[y*bounds.Dx()+x] = scaledImage.At(x, y)
		}
	}
	return averageColors, nil
}

// findImagesForSections finds the image set image that best fits the given
// chunk of sections of the target image by comparing the local color of the
// section to the average color of images in the image set.
// It processes a chunk of sections in parallel.
// Pre-compute color distances
type ColorDistance struct {
	ID       uuid.UUID
	Distance float64
}

func (cs *collageMetaDataService) findImagesForSections(
	startSection int,
	numSections int,
	sectionAverages *[]color.Color,
	imageSetAverageColors *[]*sqlc.AverageColor,
) {
	distances := make([][]ColorDistance, len(*sectionAverages))
	for i, sectionAvg := range *sectionAverages {
		distances[i] = make([]ColorDistance, len(*imageSetAverageColors))
		for j, avgColor := range *imageSetAverageColors {
			dbColor := imageprocessing.RGB16{
				R: uint16(avgColor.R) * 257,
				G: uint16(avgColor.G) * 257,
				B: uint16(avgColor.B) * 257,
			}
			distances[i][j] = ColorDistance{
				ID:       avgColor.ID,
				Distance: imageprocessing.CalculateColorDistance(dbColor, sectionAvg),
			}
		}
		// Sort once per section
		sort.Slice(distances[i], func(a, b int) bool {
			return distances[i][a].Distance < distances[i][b].Distance
		})
	}

	for section := startSection; section < startSection+numSections; section++ {
		cs.sectionMap[section] = distances[section][0].ID
	}
}

// determineImagePlacements processes the target image in batches by:
//  1. Scaling the target image to a configured resolution where the number of pixels corresponds
//     to the final collage resolution
//  2. Retrieves the average colors of each image in the image set.
//  3. Concurrently finds the best fit image by average color for each section in batches.
//  4. Encodes the placements in a metadata file stored by collage id for deferred creation.
func (cs *collageMetaDataService) determineImagePlacements() {
	cs.logger.Info(fmt.Sprintf("Finding image placements\n"))
	totalSections := cs.resolution.XSections * cs.resolution.YSections
	chunkSize := totalSections / cs.numThreads
	sectionAverages, err := cs.getSectionAverageColors()
	if err != nil {
		cs.logger.Error(fmt.Sprintf(
			"Error getting section averages\n: %s",
			err,
		))
		return
	}
	imageSetAverages, err := cs.getAverageColors()
	if err != nil {
		cs.logger.Error(fmt.Sprintf(
			"Error getting image set average colors\n: %s",
			err,
		))
		return
	}
	var wg sync.WaitGroup
	for thread := 0; thread < cs.numThreads; thread++ {
		wg.Add(1)
		threadNum := thread
		go func() {
			defer wg.Done()
			cs.findImagesForSections(
				threadNum*chunkSize,
				chunkSize,
				&sectionAverages,
				&imageSetAverages,
			)
		}()
	}
	wg.Wait()
	cs.createMetaDataFile()
}

// createMetaDataFile generates a file containing the
// image placements for the collage.
// The meta data file will be stored in the configured
// store location under the name of the collage.
func (cs *collageMetaDataService) createMetaDataFile() {
	var buf bytes.Buffer
	metaData := CollageMetaData{
		Resolution: *cs.resolution,
		SectionMap: cs.sectionMap,
	}
	err := json.NewEncoder(&buf).Encode(metaData)
	if err != nil {
		cs.logger.Error(fmt.Sprintf("Error encoding metadata file: %s\n", err))
		return
	}
	err = cs.store.PutFile(cs.collage.ID, &buf)
	if err != nil {
		cs.logger.Error(fmt.Sprintf("Error storing metadata file: %s\n", err))
	}
}
