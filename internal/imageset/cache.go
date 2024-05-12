package imageset

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/pokemonpower92/collagecommon/types"
	"github.com/pokemonpower92/imagesetservice/config"
)

// Cache is an interface that defines methods for getting and setting image sets.
type Cache interface {
	GetImageSet(key string) (*types.ImageSet, error)
	SetImageSet(imageSet *types.ImageSet) error
}

// ImageSetCache is a struct that implements the Cache interface.
type ImageSetCache struct {
	logger *log.Logger
	client *redis.Client
}

// NewImageSetCache creates a new instance of ImageSetCache.
func NewImageSetCache() *ImageSetCache {
	redisConfig := config.NewRedisConfig()
	connectionString := fmt.Sprintf(
		"%s:%s",
		redisConfig.Host,
		redisConfig.Port,
	)
	connection := redis.NewClient(&redis.Options{
		Addr:     connectionString,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	return &ImageSetCache{
		logger: log.New(log.Writer(), "cache ", log.LstdFlags),
		client: connection,
	}
}

// GetImageSet retrieves an image set from the cache based on the given key.
// It returns the image set and an error if any.
func (imageSetCache *ImageSetCache) GetImageSet(key string) (*types.ImageSet, error) {
	imageSetJson, err := imageSetCache.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	var imageSet types.ImageSet
	err = json.Unmarshal([]byte(imageSetJson), &imageSet)
	if err != nil {
		return nil, err
	}

	return &imageSet, nil
}

// SetImageSet sets the given image set in the cache.
// It returns an error if any.
func (imageSetCache *ImageSetCache) SetImageSet(imageSet *types.ImageSet) error {
	imageSetBytes, err := json.Marshal(imageSet)
	if err != nil {
		return err
	}

	key := strconv.Itoa(imageSet.ID)
	err = imageSetCache.client.Set(key, string(imageSetBytes), 0).Err()
	if err != nil {
		return err
	}

	return nil
}
