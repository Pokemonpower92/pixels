package config

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
		URI: "amqp://guest:guest@localhost:5672/",
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
	rc := RedisConfig{
		URI:      "localhost:6379",
		Password: "",
		DB:       0,
	}
	return &CacheConfig{
		RedisConfig: rc,
	}
}
