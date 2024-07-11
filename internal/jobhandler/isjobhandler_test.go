package jobhandler

import (
	"image/color"
	"log"
	"os"
	"testing"

	"github.com/pokemonpower92/collagegenerator/internal/domain"
	"github.com/pokemonpower92/collagegenerator/internal/job"
)

type mockISRepo struct{}

func (m *mockISRepo) Create(imageSet *domain.ImageSet) error {
	return nil
}

func (m *mockISRepo) Get(id int) (*domain.ImageSet, bool) {
	return nil, false
}

func (m *mockISRepo) Update(id int, imageSet *domain.ImageSet) (*domain.ImageSet, error) {
	return nil, nil
}

func (m *mockISRepo) Delete(id int) (*domain.ImageSet, error) {
	return nil, nil
}

type mockGenerator struct{}

func (m *mockGenerator) Generate(job *job.Job) (*domain.ImageSet, error) {
	return &domain.ImageSet{
		ID:          1,
		Name:        "Test",
		Description: "Test",
		AverageColors: []*color.RGBA{
			{R: 255, G: 255, B: 255, A: 255},
		},
	}, nil
}

func TestISJobHandler(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestISJobHandler",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New(os.Stdout, "test ", log.LstdFlags)
			repo := &mockISRepo{}
			generator := &mockGenerator{}
			isjh := NewISJobHandler(logger, repo, generator)

			job := &job.Job{
				ImagesetID: "1",
			}

			err := isjh.HandleJob(job)
			if err != nil {
				t.Errorf("Failed to handle job: %s", err)
			}
		})
	}
}
