package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type collageService struct {
	l                *log.Logger
	collage          *sqlc.Collage
	imageSetRepo     repository.ISRepo
	targetImageRepo  repository.TIRepo
	collageImageRepo repository.CRepo
}

func newCollageService(collage *sqlc.Collage, serviceContext context.Context) *collageService {
	postgresConfig := config.NewPostgresConfig()
	l := log.New(log.Writer(), "CollageService: ", log.LstdFlags)
	repoContext, cancel := context.WithCancel(serviceContext)
	defer cancel()

	isRepo, err := repository.NewImageSetRepository(postgresConfig, repoContext)
	if err != nil {
		l.Fatalf("Failed to create image set repository with error: %s", err)
	}
	tiRepo, err := repository.NewTagrgetImageRepository(postgresConfig, repoContext)
	if err != nil {
		l.Fatalf("Failed to create target image repository with error: %s", err)
	}
	ciRepo, err := repository.NewCollageRepository(postgresConfig, repoContext)
	if err != nil {
		l.Fatalf("Failed to create image set repository with error: %s", err)
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
func (cs *collageService) findImageForSection() {
	cs.l.Printf("Finding image for section.\n")
}

// Calculate the number of sections the target image can
// be split into given the configured resolution of the
// collage
func (cs *collageService) calculateSections() int {
	return 100
}

// Find out what image set image goes where in the collage.
func (cs *collageService) determineImagePlacements() {
	cs.l.Printf("Finding image placements.\n")
	var wg sync.WaitGroup
	numSections := cs.calculateSections()
	for i := 0; i < numSections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cs.findImageForSection()
		}()
	}
	wg.Wait()
}

func (cs *collageService) placeImagesInCollage() {
	cs.l.Printf("Placing images in the collage\n")
}
