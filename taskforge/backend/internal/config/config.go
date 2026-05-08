package config

import "os"

// Config holds application configuration
type Config struct {
	Port         string
	SecretKey    string
	TLSCert      string
	TLSKey       string
	SkipTLSVerify bool
	DevMode      bool
}

func Load() *Config {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		SecretKey:    getEnv("SECRET_KEY", "taskforge-default-secret-key-2024"),
		TLSCert:      getEnv("TLS_CERT", ""),
		TLSKey:       getEnv("TLS_KEY", ""),
		SkipTLSVerify: true,
		DevMode:      true,
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
