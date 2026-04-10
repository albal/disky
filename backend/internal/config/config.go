package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Server
	Port    string
	AppEnv  string
	AppURL  string
	JWTSecret string

	// Database
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	// Amazon PA API v5
	AmazonAccessKey  string
	AmazonSecretKey  string
	AmazonPartnerTag string
	AmazonMarketplace string
	AmazonRegion     string

	// OAuth
	GoogleClientID     string
	GoogleClientSecret string

	MicrosoftClientID     string
	MicrosoftClientSecret string
	MicrosoftTenantID     string

	AppleClientID    string
	AppleTeamID      string
	AppleKeyID       string
	ApplePrivateKey  string

	// Email / SMTP
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	c := &Config{
		Port:      getEnv("SERVER_PORT", "8080"),
		AppEnv:    getEnv("APP_ENV", "development"),
		AppURL:    getEnv("APP_URL", "http://localhost:3000"),
		JWTSecret: mustGetEnv("JWT_SECRET"),

		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),
		DBName:     getEnv("POSTGRES_DB", "disky"),
		DBUser:     getEnv("POSTGRES_USER", "disky"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "changeme"),

		AmazonAccessKey:   mustGetEnv("AMAZON_ACCESS_KEY"),
		AmazonSecretKey:   mustGetEnv("AMAZON_SECRET_KEY"),
		AmazonPartnerTag:  getEnv("AMAZON_PARTNER_TAG", "prbox-21"),
		AmazonMarketplace: getEnv("AMAZON_MARKETPLACE", "www.amazon.co.uk"),
		AmazonRegion:      getEnv("AMAZON_REGION", "eu-west-1"),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),

		MicrosoftClientID:     getEnv("MICROSOFT_CLIENT_ID", ""),
		MicrosoftClientSecret: getEnv("MICROSOFT_CLIENT_SECRET", ""),
		MicrosoftTenantID:     getEnv("MICROSOFT_TENANT_ID", "common"),

		AppleClientID:   getEnv("APPLE_CLIENT_ID", ""),
		AppleTeamID:     getEnv("APPLE_TEAM_ID", ""),
		AppleKeyID:      getEnv("APPLE_KEY_ID", ""),
		ApplePrivateKey: loadAppleKey(),

		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     smtpPort,
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "Disky <noreply@disky.tsew.com>"),
	}

	return c, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPassword)
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		// In non-production / test, allow empty values
		return ""
	}
	return v
}

func loadAppleKey() string {
	path := os.Getenv("APPLE_PRIVATE_KEY_PATH")
	if path == "" {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
