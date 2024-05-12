package imageset

import (
	"image/color"
	"log"
	"testing"

	"github.com/pokemonpower92/collagecommon/types"
)

type mockLogger struct{}

func (ml *mockLogger) Printf(format string, v ...interface{}) {}

type mockCache struct {
	expectedImageSet *types.ImageSet
	expectedErr      error
}

func (mc *mockCache) GetImageSet(key string) (*types.ImageSet, error) {
	return mc.expectedImageSet, mc.expectedErr
}

func (mc *mockCache) SetImageSet(im *types.ImageSet) error {
	return mc.expectedErr
}

type mockDB struct {
	expectedImageSet *types.ImageSet
	expectedErr      error
}

func (mdb *mockDB) GetImageSet(id int) (*types.ImageSet, error) {
	return mdb.expectedImageSet, mdb.expectedErr
}

func (mdb *mockDB) CreateImageSet(im *types.ImageSet) error {
	return mdb.expectedErr
}

func (mdb *mockDB) SetAverageColors(id int, ac []*color.RGBA) error {
	return mdb.expectedErr
}

func TestHandleJob(t *testing.T) {
	subtests := []struct {
		name          string
		job           *Job
		expectedErr   error
		expectedImage *types.ImageSet
	}{
		{
			name: "happy path",
			job: &Job{
				ImagesetID:  "123",
				BucketName:  "bucket",
				Description: "description",
			},
			expectedErr: nil,
			expectedImage: &types.ImageSet{
				ID:            123,
				Name:          "bucket",
				Description:   "description",
				AverageColors: []*color.RGBA{{R: 127, G: 127, B: 127, A: 255}},
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			jh := &JobHandler{
				logger: log.New(log.Writer(), "jobhandler ", log.LstdFlags),
				cache: &mockCache{
					expectedImageSet: subtest.expectedImage,
					expectedErr:      subtest.expectedErr,
				},
				db: &mockDB{
					expectedImageSet: subtest.expectedImage,
					expectedErr:      subtest.expectedErr,
				},
			}

			jh.HandleJob(subtest.job)
		})
	}
}
