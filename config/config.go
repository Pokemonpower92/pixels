package config

import (
	"os"
)

func PrivateKeyPem() string {
	return os.Getenv("PRIVATE_KEY_PEM")
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	Port     string
	DBName   string
}

func NewPostgresConfig() *DBConfig {
	return &DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}
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
