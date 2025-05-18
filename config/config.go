package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Addr         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LogLevel     string
	Db           DbConfig
	TokenConfig  TokenConfig
	Env          string
	APIURL       string
	FrontendURL  string
	Auth         AuthConfig
	Redis        RedisConfig
	RateLimiter  RateLimiterConfig
}

// All configuration structs now use exported fields

type RedisConfig struct {
	Addr    string
	Pw      string
	DB      int
	Enabled bool
}

type AuthConfig struct {
	Basic BasicConfig
	Token TokenConfig
}

type TokenConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
}

type BasicConfig struct {
	User string
	Pass string
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
	Username     string
	Password     string
	Database     string
	Host         string
}

type RateLimiterConfig struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	var config Config

	// Server config
	config.Addr = envOrDefault("PORT", "8080")
	if !hasPrefix(config.Addr, ":") {
		config.Addr = ":" + config.Addr
	}

	config.IdleTimeout = parseDuration(envOrDefault("IDLE_TIMEOUT", "60s"))
	config.ReadTimeout = parseDuration(envOrDefault("READ_TIMEOUT", "10s"))
	config.WriteTimeout = parseDuration(envOrDefault("WRITE_TIMEOUT", "30s"))
	config.LogLevel = envOrDefault("LOG_LEVEL", "info")
	config.Env = envOrDefault("APP_ENV", "development")
	config.APIURL = envOrDefault("API_URL", "localhost:8080")
	config.FrontendURL = envOrDefault("FRONTEND_URL", "http://localhost:5173")

	// JWT Token config
	config.TokenConfig.Secret = envOrDefault("JWT_SECRET", "")
	if config.TokenConfig.Secret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET environment variable must be set")
	}
	config.TokenConfig.Exp = parseDuration(envOrDefault("TOKEN_EXP", "24h"))
	config.TokenConfig.Iss = envOrDefault("TOKEN_ISS", "JonoMot")

	// Database config
	config.Db.Username = envOrDefault("DB_USERNAME", "admin")
	config.Db.Password = envOrDefault("DB_PASSWORD", "adminpassword")
	config.Db.Database = envOrDefault("DB_DATABASE", "jono-poll-db")
	config.Db.Host = envOrDefault("DB_HOST", "localhost") // Add this to allow overriding in Docker

	// Build the connection string programmatically
	config.Db.Addr = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.Db.Username,
		config.Db.Password,
		config.Db.Host,
		config.Db.Database)

	config.Db.MaxOpenConns = parseInt(envOrDefault("DB_MAX_OPEN_CONNS", "30"))
	config.Db.MaxIdleConns = parseInt(envOrDefault("DB_MAX_IDLE_CONNS", "30"))
	config.Db.MaxIdleTime = envOrDefault("DB_MAX_IDLE_TIME", "15m")

	// Rate limiter config
	config.RateLimiter.RequestsPerTimeFrame = parseInt(envOrDefault("RATELIMITER_REQUESTSPERTIMEFRAME", "20"))
	config.RateLimiter.TimeFrame = parseDuration(envOrDefault("RATELIMITER_TIMEFRAME", "5s"))
	config.RateLimiter.Enabled = parseBool(envOrDefault("RATELIMITER_ENABLED", "true"))

	return config, nil
}

// Helper functions for env variable parsing

// envOrDefault gets an environment variable or returns the default if not set
func envOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// parseDuration parses a duration from string with a fallback to a default
func parseDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		// Log error or handle it as appropriate
		fmt.Printf("Warning: invalid duration '%s', using 0\n", value)
		return 0
	}
	return duration
}

// parseInt parses an integer with a fallback to 0
func parseInt(value string) int {
	intVal, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("Warning: invalid integer '%s', using 0\n", value)
		return 0
	}
	return intVal
}

// parseBool parses a boolean with a fallback
func parseBool(value string) bool {
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		fmt.Printf("Warning: invalid boolean '%s', using false\n", value)
		return false
	}
	return boolVal
}

// hasPrefix checks if a string has a certain prefix
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}
