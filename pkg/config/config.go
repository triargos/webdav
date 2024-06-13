package config

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/triargos/webdav/pkg/logging"
	"os"
)

type Config struct {
	Network *NetworkConfig   `mapstructure:"network"`
	Content *ContentConfig   `mapstructure:"content"`
	Users   *map[string]User `mapstructure:"users"`
}

type NetworkConfig struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
	Prefix  string `mapstructure:"prefix,omitempty"`
}

type ContentConfig struct {
	Dir string `mapstructure:"dir"`
}

type User struct {
	Password       string   `mapstructure:"password"`
	Root           string   `mapstructure:"root,omitempty"`
	SubDirectories []string `mapstructure:"subdirectories,omitempty"`
	Jail           bool     `mapstructure:"jail,omitempty"`
	Admin          bool     `mapstructure:"admin"`
}

func getConfigurationPath() string {
	configPath := "/etc/webdav"
	if os.Getenv("DOCKER_ENABLED") != "1" {
		configPath = "./config"
	}
	return configPath

}

var defaultConfig = Config{
	Network: &NetworkConfig{
		Address: "0.0.0.0",
		Port:    "8080",
		Prefix:  "/",
	},
	Content: &ContentConfig{
		Dir: "/var/webdav/data",
	},
	Users: &map[string]User{
		"admin": {
			Password:       "admin",
			Admin:          true,
			Jail:           false,
			Root:           "/Users/admin",
			SubDirectories: []string{"documents"},
		},
	},
}

func Get() *Config {
	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return &defaultConfig
	}
	return cfg
}

func Set(cfg *Config) {
	viper.Set("network", cfg.Network)
	viper.Set("content", cfg.Content)
	viper.Set("users", cfg.Users)
}

func Read() error {
	path := getConfigurationPath()
	logging.Log.Info.Printf("Reading configuration from %s\n", path)
	viper.AddConfigPath(getConfigurationPath())
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			logging.Log.Info.Println("No configuration file found, creating default configuration")
			if err := WriteDefaultConfig(); err != nil {
				return err
			}
		}
	}
	return nil
}

func WriteDefaultConfig() error {
	viper.AddConfigPath(getConfigurationPath())
	viper.Set("network", defaultConfig.Network)
	viper.Set("content", defaultConfig.Content)
	viper.Set("users", defaultConfig.Users)
	return viper.SafeWriteConfig()
}

func AddUser(username string, user User) {
	cfg := Get()
	users := *cfg.Users
	users[username] = user
	cfg.Users = &users
	Set(cfg)
}

func Write() error {
	return viper.WriteConfig()
}
