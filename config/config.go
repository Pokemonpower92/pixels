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
	Host     string
	Port     string
	User     string
	Password string
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
		Host:     os.Getenv("RABBITMQ_HOST"),
		Port:     os.Getenv("RABBITMQ_PORT"),
		User:     os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
	}

	return &ConsumerConfig{
		Queue:          qc,
		RabbitMQConfig: rmqc,
	}
}

type RedisConfig struct {
	Host     string
	User     string
	Password string
	Port     string
	DB       int
}

func NewRedisConfig() *RedisConfig {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic("Invalid REDIS_DB value")
	}

	return &RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		User:     os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASSWORD"),
		Port:     os.Getenv("REDIS_PORT"),
		DB:       db,
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

type DBConfig struct {
	Host     string
	User     string
	Password string
	Port     string
	DbName   string
}

func NewDBConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DbName:   os.Getenv("IMAGESET_DB"),
	}
}
