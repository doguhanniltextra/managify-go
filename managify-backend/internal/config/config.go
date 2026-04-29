package config

import (
	"fmt"
	"os"
)

type Config struct {
	SecretKey   string
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	SMTPFrom    string
	FrontendURL string
}

func LoadConfig() *Config {
	return &Config{
		SecretKey:   os.Getenv("SECRET_KEY"),
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    os.Getenv("SMTP_PORT"),
		SMTPUser:    os.Getenv("SMTP_USER"),
		SMTPPass:    os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:    os.Getenv("SMTP_FROM"),
		FrontendURL: os.Getenv("FRONTEND_URL"),
	}
}

func (c *Config) Validate() error {
	if c.SecretKey == "" {
		return fmt.Errorf("SECRET_KEY environment variable is not set or empty")
	}
	return nil
}
