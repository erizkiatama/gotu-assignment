package config

import (
	"github.com/spf13/viper"
)

var config *Config

// option defines configuration options
type option struct {
	configFolder string
	configFile   string
	configType   string
}

// Option define an option for config package
type Option func(*option)

// WithConfigFolder set `config` to use the given config folder
func WithConfigFolder(configFolder string) Option {
	return func(opt *option) {
		opt.configFolder = configFolder
	}
}

// WithConfigFile set `config` to use the given config file
func WithConfigFile(configFile string) Option {
	return func(opt *option) {
		opt.configFile = configFile
	}
}

// WithConfigType set `config` to use the given config type
func WithConfigType(configType string) Option {
	return func(opt *option) {
		opt.configType = configType
	}
}

// getDefaultConfigFolder get default config folder.
func getDefaultConfigFolder() string {
	return "./internal/config/"
}

// getDefaultConfigFile get default config file.
func getDefaultConfigFile() string {
	return "config"
}

// getDefaultConfigType get default config type.
func getDefaultConfigType() string {
	return "yaml"
}

// Get config
func Get() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}

func Load(opts ...Option) (*Config, error) {
	opt := &option{
		configFolder: getDefaultConfigFolder(),
		configFile:   getDefaultConfigFile(),
		configType:   getDefaultConfigType(),
	}

	for _, optFunc := range opts {
		optFunc(opt)
	}

	// Config File Path
	viper.AddConfigPath(opt.configFolder)
	// Config File Name
	viper.SetConfigName(opt.configFile)
	// Config File Type
	viper.SetConfigType(opt.configType)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config = new(Config)
	if err = viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Reading secret variables that will injected in production settings
	if config.Server.Env == "production" {
		config.Database.Postgres.Host = viper.GetString("APP_POSTGRES_HOST")
		config.Database.Postgres.Username = viper.GetString("APP_POSTGRES_USERNAME")
		config.Database.Postgres.Password = viper.GetString("APP_POSTGRES_PASSWORD")
		config.Database.Postgres.Database = viper.GetString("APP_POSTGRES_DATABASE")
		config.Server.SecretKey = viper.GetString("APP_SERVER_SECRET_KEY")
	}

	return config, nil
}
