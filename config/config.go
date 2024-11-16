package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"strings"
	"time"
)

type ConfigStruct struct {
	Addr        string
	LogLevel    string
	db          dbConfigStruct
	Env         string
	apiURL      string
	mail        mailConfigStruct
	frontendURL string
	auth        authConfigStruct
	redisCfg    redisConfigStruct
	rateLimiter rateLimiterConfigStruct
}

type redisConfigStruct struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

type authConfigStruct struct {
	basic basicConfigStruct
	token tokenConfigStruct
}

type tokenConfigStruct struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfigStruct struct {
	user string
	pass string
}

type mailConfigStruct struct {
	sendGrid  sendGridConfigStruct
	fromEmail string
	exp       time.Duration
}

type sendGridConfigStruct struct {
	apiKey string
}

type dbConfigStruct struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type rateLimiterConfigStruct struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}

// LoadConfig loads the configuration from the specified filename and environment variables.
func LoadConfig(filename string) (ConfigStruct, error) {
	// Create a new Viper instance.
	v := viper.New()

	// Set the configuration file name and path.
	v.SetConfigName(filename)
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	// Enable reading from environment variables.
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the configuration file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist, we'll use env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Error reading config file: %v\n", err)
			return ConfigStruct{}, err
		}
	}

	// Automatically bind environment variables based on struct tags
	if err := bindEnvVariables(v); err != nil {
		return ConfigStruct{}, fmt.Errorf("error binding environment variables: %v", err)
	}
	// Add more bindings for other config fields as needed

	// Set default values
	setDefaultValues(v)

	// Unmarshal the configuration into the Config struct.
	var config ConfigStruct
	if err := v.Unmarshal(&config); err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		return ConfigStruct{}, err
	}

	return config, nil
}

// bindEnvVariables automatically binds environment variables based on struct tags
func bindEnvVariables(v *viper.Viper) error {
	configType := reflect.TypeOf(ConfigStruct{})
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		envTag := field.Tag.Get("env")
		if envTag != "" {
			err := v.BindEnv(field.Name, envTag)
			if err != nil {
				return fmt.Errorf("error binding %s: %v", field.Name, err)
			}
		}
	}
	return nil
}

// setDefaultValues sets default values for configuration fields
func setDefaultValues(v *viper.Viper) {
	v.SetDefault("Addr", ":8080")
	v.SetDefault("LogLevel", "info")
	// Add more default values as needed
}
