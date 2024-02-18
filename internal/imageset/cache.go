package imageset

import (
	"encoding/json"
	"log"

	"github.com/go-redis/redis"
	"github.com/pokemonpower92/imagesetservice/config"
)

type Cache struct {
	l    *log.Logger
	conn *redis.Client
}

func NewCache(l *log.Logger) *Cache {
	cc := config.NewCacheConfig()
	conn := redis.NewClient(&redis.Options{
		Addr:     cc.RedisConfig.URI,
		Password: cc.RedisConfig.Password,
		DB:       cc.RedisConfig.DB,
	})

	return &Cache{
		l:    l,
		conn: conn,
	}
}

func (c *Cache) GetImageSet(id string) (*ImageSet, error) {
	val, err := c.conn.Get(id).Result()
	if err != nil {
		return nil, err
	}
	var im ImageSet
	err = json.Unmarshal([]byte(val), &im)
	if err != nil {
		return nil, err
	}

	return &im, nil
}

func (c *Cache) SetImageSet(im *ImageSet) error {
	b, err := json.Marshal(im)
	if err != nil {
		return err
	}
	err = c.conn.Set(im.ID, string(b), 0).Err()
	if err != nil {
		return err
	}

	return nil
}
