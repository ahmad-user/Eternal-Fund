package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type DbConfig struct {
	Host       string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	Driver     string
}

type ApiConfig struct {
	ApiPort string
}

type TokenConfig struct {
	IssuerName    string
	SignatureKey  []byte
	SigningMethod *jwt.SigningMethodHMAC
	ExpiresTime   time.Duration
}

type Config struct {
	DbConfig
	ApiConfig
	TokenConfig
}

func (c *Config) Configuration() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("missing env file %v", err.Error())
	}
	c.DbConfig = DbConfig{
		Host:       os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		Driver:     os.Getenv("DB_DRIVER"),
	}

	c.ApiConfig = ApiConfig{ApiPort: os.Getenv("API_PORT")}

	tokenExpire, _ := strconv.Atoi(os.Getenv("TOKEN_EXPIRE"))

	c.TokenConfig = TokenConfig{
		IssuerName:    os.Getenv("TOKEN_ISSUE"),
		SignatureKey:  []byte(os.Getenv("TOKEN_SECRET")),
		SigningMethod: jwt.SigningMethodHS256,
		ExpiresTime:   time.Duration(tokenExpire) * time.Minute,
	}

	if c.Host == "" || c.DbPort == "" || c.DbUser == "" ||
		c.DbPassword == "" || c.DbName == "" || c.Driver == "" || c.IssuerName == "" ||
		len(c.SignatureKey) == 0 || c.ExpiresTime < 0 {
		return fmt.Errorf("missing enivronment")
	}
	return nil
}

func NewConfig() (*Config, error) {
	config := &Config{}

	if err := config.Configuration(); err != nil {
		return nil, err
	}
	return config, nil
}
