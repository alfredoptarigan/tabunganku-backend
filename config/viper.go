package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"

	"alfredo/tabunganku/pkg/log"
)

// DatabaseConfig holds all database configuration
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	SSLMode         string `mapstructure:"sslmode"`
	Timezone        string `mapstructure:"timezone"`
	MaxConnections  int    `mapstructure:"max_connections"`
	IdleConnections int    `mapstructure:"max_idle_connections"`
	MaxLifetime     int    `mapstructure:"conn_max_lifetime"`
}

// ApplicationConfig holds the application configuration
type ApplicationConfig struct {
	Port        string `mapstructure:"port"`
	Environment string `mapstructure:"environment"`
}

// RedisConfig holds all Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// Config is the main configuration structure
type Config struct {
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	Application ApplicationConfig `mapstructure:"application"`
}

// ViperConfig is an interface for the Viper configuration library
type ViperConfig interface {
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetStringSlice(key string) []string
	GetIntSlice(key string) []int
	GetTime(key string) time.Time
	UnmarshalKey(key string, rawVal interface{}) error
}

// viperWrapper is a wrapper around the Viper library
type viperWrapper struct {
	viper *viper.Viper
}

// NewViperConfig initializes a new Viper configuration
// Dalam function NewViperConfig(), tambahkan:
func NewViperConfig() ViperConfig {
	v := viper.New()

	// Check if running in test environment
	if os.Getenv("GO_ENV") == "test" {
		v.SetConfigName("config.test")
	} else {
		v.SetConfigName("config")
	}

	v.SetConfigType("yaml")
	v.AddConfigPath(setRootPath())

	v.AutomaticEnv()

	// read the config file
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			panic("config file not found")
		}
		panic("failed to read config file: " + err.Error())
	}

	// unmarshal the config file into the Config struct
	if err := v.Unmarshal(&v); err != nil {
		panic("failed to unmarshal config file")
	}

	return &viperWrapper{viper: v}
}

func GetLoggingConfig() log.LoggingConfig {
	v := NewViperConfig()
	var loggingConfig log.LoggingConfig

	// Unmarshal the logging configuration
	if err := v.UnmarshalKey("logging", &loggingConfig); err != nil {
		panic("failed to unmarshal logging config: " + err.Error())
	}

	return loggingConfig
}

// Get retrieves a value from the Viper configuration
func (v *viperWrapper) Get(key string) interface{} {
	return v.viper.Get(key)
}

// GetString retrieves a string value from the Viper configuration
func (v *viperWrapper) GetString(key string) string {
	return v.viper.GetString(key)
}

// GetInt retrieves an integer value from the Viper configuration
func (v *viperWrapper) GetInt(key string) int {
	return v.viper.GetInt(key)
}

// GetBool retrieves a boolean value from the Viper configuration
func (v *viperWrapper) GetBool(key string) bool {
	return v.viper.GetBool(key)
}

// GetFloat64 retrieves a float64 value from the Viper configuration
func (v *viperWrapper) GetFloat64(key string) float64 {
	return v.viper.GetFloat64(key)
}

// GetStringSlice retrieves a slice of strings from the Viper configuration
func (v *viperWrapper) GetStringSlice(key string) []string {
	return v.viper.GetStringSlice(key)
}

// GetIntSlice retrieves a slice of integers from the Viper configuration
func (v *viperWrapper) GetIntSlice(key string) []int {
	return v.viper.GetIntSlice(key)
}

// GetTime retrieves a time.Time value from the Viper configuration
func (v *viperWrapper) GetTime(key string) time.Time {
	return v.viper.GetTime(key)
}

// UnmarshalKey unmarshals a key from the Viper configuration into a struct
func (v *viperWrapper) UnmarshalKey(key string, rawVal interface{}) error {
	return v.viper.UnmarshalKey(key, rawVal)
}

func setRootPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory: " + err.Error())
	}

	if strings.Contains(wd, "/cmd/cron") {
		return filepath.Join(wd, "../../")
	}
	if strings.Contains(wd, "/tests/integration") {
		return filepath.Join(wd, "../../")
	}

	return wd
}
