package config

import (
	"os"
)

type Config struct {
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	SMTPFrom    string
	FrontendURL string
}

func LoadConfig() *Config {
	return &Config{
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    os.Getenv("SMTP_PORT"),
		SMTPUser:    os.Getenv("SMTP_USER"),
		SMTPPass:    os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:    os.Getenv("SMTP_FROM"),
		FrontendURL: os.Getenv("FRONTEND_URL"),
	}
}
