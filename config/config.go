package config

import (
	"os"
	"strconv"
	"sync"
)

var (
	once sync.Once
	env  map[string]string
)

type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       map[string]interface{}
}

type RabbitMQConfig struct {
	URI string
}

type ConsumerConfig struct {
	Queue          QueueConfig
	RabbitMQConfig RabbitMQConfig
}

func NewConsumerConfig() *ConsumerConfig {
	qc := QueueConfig{
		Name:       "hello",
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
	rmqc := RabbitMQConfig{
		URI: os.Getenv("RABBITMQ_URI"),
	}

	return &ConsumerConfig{
		Queue:          qc,
		RabbitMQConfig: rmqc,
	}
}

type RedisConfig struct {
	URI      string
	Password string
	DB       int
}

type CacheConfig struct {
	RedisConfig RedisConfig
}

func NewCacheConfig() *CacheConfig {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic("Invalid REDIS_DB value")
	}

	rc := RedisConfig{
		URI:      os.Getenv("REDIS_URI"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}

	return &CacheConfig{
		RedisConfig: rc,
	}
}

type S3Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

func NewS3Config() *S3Config {
	return &S3Config{
		Region:          os.Getenv("S3_REGION"),
		AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
	}
}
