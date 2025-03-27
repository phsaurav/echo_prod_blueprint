package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"strings"
	"time"
)

var (
	QueryTimeoutDuration = time.Second * 5
)

type Config struct {
	Addr         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LogLevel     string
	Db           dbConfig
	Env          string
	apiURL       string
	mail         mailConfig
	frontendURL  string
	auth         authConfig
	redisCfg     redisConfig
	rateLimiter  rateLimiterConfig
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	user string
	pass string
}

type mailConfig struct {
	sendGrid  sendGridConfig
	fromEmail string
	exp       time.Duration
}

type sendGridConfig struct {
	apiKey string
}

type dbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type rateLimiterConfig struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}

// LoadConfig loads the configuration from the specified filename and environment variables.
func LoadConfig(filename string) (Config, error) {
	// Create a new Viper instance.
	cfgReader := viper.New()

	// Set the configuration file name and path.
	cfgReader.SetConfigName(filename)
	cfgReader.AddConfigPath("./config")
	cfgReader.AddConfigPath(".")

	// Enable reading from environment variables.
	cfgReader.AutomaticEnv()
	cfgReader.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicit bindings for DB settings
	cfgReader.BindEnv("Db.Addr", "DB_ADDR")
	cfgReader.BindEnv("Db.MaxIdleTime", "DB_MAX_IDLE_TIME")
	cfgReader.BindEnv("Db.MaxOpenConns", "DB_MAX_OPEN_CONNS")
	cfgReader.BindEnv("Db.MaxIdleConns", "DB_MAX_IDLE_CONNS")

	// Read the configuration file.
	if err := cfgReader.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist, we'll use env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Error reading config file: %cfgReader\n", err)
			return Config{}, err
		}
	}

	// Automatically bind environment variables based on struct tags
	if err := bindEnvVariables(cfgReader); err != nil {
		return Config{}, fmt.Errorf("error binding environment variables: %cfgReader", err)
	}
	// Add more bindings for other config fields as needed

	// Set default values
	setDefaultValues(cfgReader)

	// Unmarshal the configuration into the Config struct.
	var config Config
	if err := cfgReader.Unmarshal(&config); err != nil {
		fmt.Printf("Error unmarshaling config: %cfgReader\n", err)
		return Config{}, err
	}

	return config, nil
}

// bindEnvVariables automatically binds environment variables based on struct tags
func bindEnvVariables(cfgReader *viper.Viper) error {
	configType := reflect.TypeOf(Config{})
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		envTag := field.Tag.Get("env")
		if envTag != "" {
			err := cfgReader.BindEnv(field.Name, envTag)
			if err != nil {
				return fmt.Errorf("error binding %s: %v", field.Name, err)
			}
		}
	}
	return nil
}

// setDefaultValues sets default values for configuration fields
func setDefaultValues(cfgReader *viper.Viper) {
	cfgReader.SetDefault("Addr", ":8080")
	cfgReader.SetDefault("LogLevel", "info")
	cfgReader.SetDefault("IdleTimeout", time.Minute)
	cfgReader.SetDefault("ReadTimeout", 10*time.Second)
	cfgReader.SetDefault("WriteTimeout", 30*time.Second)
	cfgReader.SetDefault("db.MaxIdleTime", "15m")
	// Add more default values as needed
}
