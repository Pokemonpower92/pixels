package jobhandler

import (
	"errors"
	"log"
	"strconv"
	"testing"

	"github.com/pokemonpower92/collagegenerator/internal/domain"
	"github.com/pokemonpower92/collagegenerator/internal/job"
	"github.com/pokemonpower92/collagegenerator/internal/stubs"
)

func TestISJobHandler(t *testing.T) {
	tests := []struct {
		name         string
		job          *job.ImageSetJob
		getFunc      func(id int) (*domain.ImageSet, bool)
		createFunc   func(obj *domain.ImageSet) error
		generateFunc func(job *job.ImageSetJob) (*domain.ImageSet, error)
		expected     error
	}{
		{
			name: "invalid imageset id",
			job: &job.ImageSetJob{
				ImagesetID:  "bad id",
				Path:        "test_bucket",
				Description: "test description",
			},
			getFunc: func(id int) (*domain.ImageSet, bool) {
				return nil, false
			},
			createFunc: func(obj *domain.ImageSet) error {
				return nil
			},
			generateFunc: func(job *job.ImageSetJob) (*domain.ImageSet, error) {
				return nil, nil
			},
			expected: &strconv.NumError{Func: "Atoi", Num: "bad id", Err: strconv.ErrSyntax},
		},
		{
			name: "imageset found",
			job: &job.ImageSetJob{
				ImagesetID:  "1",
				Path:        "test_bucket",
				Description: "test description",
			},
			getFunc: func(id int) (*domain.ImageSet, bool) {
				return &domain.ImageSet{}, true
			},
			createFunc: func(obj *domain.ImageSet) error {
				return nil
			},
			generateFunc: func(job *job.ImageSetJob) (*domain.ImageSet, error) {
				return &domain.ImageSet{}, nil
			},
			expected: errors.New(""),
		},
		{
			name: "imageset not found",
			job: &job.ImageSetJob{
				ImagesetID:  "1",
				Path:        "test_bucket",
				Description: "test description",
			},
			getFunc: func(id int) (*domain.ImageSet, bool) {
				return nil, false
			},
			createFunc: func(obj *domain.ImageSet) error {
				return nil
			},
			generateFunc: func(job *job.ImageSetJob) (*domain.ImageSet, error) {
				return &domain.ImageSet{}, nil
			},
			expected: errors.New(""),
		},
		{
			name: "failed to generate imageset",
			job: &job.ImageSetJob{
				ImagesetID:  "1",
				Path:        "test_bucket",
				Description: "test description",
			},
			getFunc: func(id int) (*domain.ImageSet, bool) {
				return nil, false
			},
			createFunc: func(obj *domain.ImageSet) error {
				return nil
			},
			generateFunc: func(job *job.ImageSetJob) (*domain.ImageSet, error) {
				return nil, errors.New("failed to generate imageset")
			},
			expected: errors.New("failed to generate imageset"),
		},
		{
			name: "failed to create imageset",
			job: &job.ImageSetJob{
				ImagesetID:  "1",
				Path:        "test_bucket",
				Description: "test description",
			},
			getFunc: func(id int) (*domain.ImageSet, bool) {
				return nil, false
			},
			createFunc: func(obj *domain.ImageSet) error {
				return errors.New("failed to create imageset")
			},
			generateFunc: func(job *job.ImageSetJob) (*domain.ImageSet, error) {
				return &domain.ImageSet{}, nil
			},
			expected: errors.New("failed to create imageset"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubs.RepositoryStub[domain.ImageSet]{
				GetFunc:    tt.getFunc,
				CreateFunc: tt.createFunc,
			}

			gen := &stubs.GeneratorStub{
				GenerateFunc: tt.generateFunc,
			}

			logger := log.New(log.Writer(), "test: ", log.Flags())
			handler := NewISJobHandler(logger, repo, gen)

			err := handler.HandleJob(tt.job)
			if errors.Is(err, tt.expected) {
				t.Errorf("expected error %v, got %v", tt.expected, err)
			}
		})
	}
}
