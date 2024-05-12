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

type iCache interface {
	GetImageSet(key string) (*types.ImageSet, error)
	SetImageSet(im *types.ImageSet) error
}

type Cache struct {
	l    *log.Logger
	conn *redis.Client
}

func NewCache() *Cache {
	cc := config.NewCacheConfig()
	connString := fmt.Sprintf("%s:%s", cc.RedisConfig.Host, cc.RedisConfig.Port)
	conn := redis.NewClient(&redis.Options{
		Addr:     connString,
		Password: cc.RedisConfig.Password,
		DB:       cc.RedisConfig.DB,
	})

	return &Cache{
		l:    log.New(log.Writer(), "cache ", log.LstdFlags),
		conn: conn,
	}
}

func (c *Cache) GetImageSet(key string) (*types.ImageSet, error) {
	val, err := c.conn.Get(key).Result()
	if err != nil {
		return nil, err
	}
	var im types.ImageSet
	err = json.Unmarshal([]byte(val), &im)
	if err != nil {
		return nil, err
	}

	return &im, nil
}

func (c *Cache) SetImageSet(im *types.ImageSet) error {
	b, err := json.Marshal(im)
	if err != nil {
		return err
	}

	key := strconv.Itoa(im.ID)
	err = c.conn.Set(key, string(b), 0).Err()
	if err != nil {
		return err
	}

	return nil
}
