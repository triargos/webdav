package config

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/triargos/webdav/pkg/logging"
	"os"
	"path/filepath"
)

var v = viper.NewWithOptions(
	viper.KeyDelimiter("::"))

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
	if err := v.Unmarshal(&cfg); err != nil {
		return &defaultConfig
	}
	return cfg
}

func Set(cfg *Config) {
	v.Set("network", cfg.Network)
	v.Set("content", cfg.Content)
	v.Set("users", cfg.Users)
}

func Read() error {
	path := getConfigurationPath()
	logging.Log.Info.Printf("Reading configuration from %s\n", path)
	v.SetConfigFile(filepath.Join(path, "config.yaml"))
	//Set default values
	v.SetDefault("network", defaultConfig.Network)
	v.SetDefault("content", defaultConfig.Content)
	v.SetDefault("users", defaultConfig.Users)
	v.SetConfigType("yaml")
	readErr := v.ReadInConfig()
	var configFileNotFoundError viper.ConfigFileNotFoundError
	if errors.As(readErr, &configFileNotFoundError) {
		logging.Log.Info.Println("No configuration file found, creating default configuration")
		v.WriteConfig()
	}
	return nil
}

func WriteDefaultConfig() error {
	v.AddConfigPath(getConfigurationPath())
	v.Set("network", defaultConfig.Network)
	v.Set("content", defaultConfig.Content)
	v.Set("users", defaultConfig.Users)
	return v.SafeWriteConfig()
}

func AddUser(username string, user User) {
	cfg := Get()
	users := *cfg.Users
	users[username] = user
	cfg.Users = &users
	Set(cfg)
}

func RemoveUser(username string) {
	cfg := Get()
	users := *cfg.Users
	delete(users, username)
	cfg.Users = &users
	Set(cfg)
}

func Write() error {
	return v.WriteConfig()
}
