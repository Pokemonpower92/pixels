package service

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type collageService struct {
	l                *log.Logger
	collage          *sqlc.Collage
	imageSetRepo     repository.ISRepo
	targetImageRepo  repository.TIRepo
	collageImageRepo repository.CRepo
	store            datastore.Store
}

func newCollageService(collage *sqlc.Collage, serviceContext context.Context) *collageService {
	postgresConfig := config.NewPostgresConfig()
	l := log.New(log.Writer(), "CollageService: ", log.LstdFlags)
	isRepo, err := repository.NewImageSetRepository(postgresConfig, serviceContext)
	if err != nil {
		l.Fatalf("Failed to create image set repository with error: %s", err)
	}
	tiRepo, err := repository.NewTagrgetImageRepository(postgresConfig, serviceContext)
	if err != nil {
		l.Fatalf("Failed to create target image repository with error: %s", err)
	}
	ciRepo, err := repository.NewCollageRepository(postgresConfig, serviceContext)
	if err != nil {
		l.Fatalf("Failed to create collage repository with error: %s", err)
	}
	return &collageService{
		l:                l,
		collage:          collage,
		imageSetRepo:     isRepo,
		targetImageRepo:  tiRepo,
		collageImageRepo: ciRepo,
	}
}

func CreateCollage(collage *sqlc.Collage) {
	serviceContext, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	service := newCollageService(collage, serviceContext)
	service.determineImagePlacements()
	service.placeImagesInCollage()
}

// Find the image set image that best fits the given
// section of the target image
func (cs *collageService) findImageForSection(section int) {
	cs.l.Printf("Finding image for section: %d.\n", section)
}

// Calculate the number of sections the target image can
// be split into given the configured resolution of the
// collage
func (cs *collageService) calculateSections() int {
	targetImage, err := cs.targetImageRepo.Get(cs.collage.TargetImageID)
	if err != nil {
		cs.l.Fatalf("Failed to get target image with error: %s\n", err)
	}
	return len(strings.Split(targetImage.Name, ""))
}

// Find out what image set image goes where in the collage.
func (cs *collageService) determineImagePlacements() {
	cs.l.Printf("Finding image placements.\n")
	var wg sync.WaitGroup
	numSections := cs.calculateSections()
	for section := 0; section < numSections; section++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cs.findImageForSection(section)
		}()
	}
	wg.Wait()
}

func (cs *collageService) placeImagesInCollage() {
	cs.l.Printf("Placing images in the collage\n")
}
