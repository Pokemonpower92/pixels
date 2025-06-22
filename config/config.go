package config

import (
	"fmt"
	"os"
)

func PrivateKeyPem() string {
	return os.Getenv("PRIVATE_KEY_PEM")
}

func ConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
}

type ResolutionConfig struct {
	CollageWidth  int
	CollageHeight int
	SectionWidth  int
	SectionHeight int
	XSections     int
	YSections     int
}

func NewResolutionConfig() *ResolutionConfig {
	return &ResolutionConfig{
		CollageWidth:  64,
		CollageHeight: 64,
		SectionWidth:  1,
		SectionHeight: 1,
		XSections:     64,
		YSections:     64,
	}
}

type ServerConfig struct {
	Host string
	Port string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Host: os.Getenv("SERVER_HOST"),
		Port: os.Getenv("SERVER_PORT"),
	}
}
