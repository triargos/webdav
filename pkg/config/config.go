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

var configTemplate = Config{
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

func DeepCopyConfig(original Config) Config {
	newConfig := Config{
		Network: &NetworkConfig{
			Address: original.Network.Address,
			Port:    original.Network.Port,
			Prefix:  original.Network.Prefix,
		},
		Content: &ContentConfig{
			Dir: original.Content.Dir,
		},
		Users: &map[string]User{},
	}

	for k, v := range *original.Users {
		newUser := v
		(*newConfig.Users)[k] = newUser
	}

	return newConfig
}

func Get() *Config {
	cfg := &Config{}
	if err := v.Unmarshal(&cfg); err != nil {
		webdavPort, createNoAdminUser, webdavDataDir := ReadEnv()
		defaultConfig := GenerateDefaultConfig(webdavPort, createNoAdminUser, webdavDataDir)
		return &defaultConfig
	}
	return cfg
}

func Set(cfg *Config) {
	v.Set("network", cfg.Network)
	v.Set("content", cfg.Content)
	v.Set("users", cfg.Users)
}

func initViper() {
	v.SetConfigType("yaml")
}

func ReadEnv() (string, bool, string) {
	webdavPort := os.Getenv("WEBDAV_PORT")
	createNoAdminUser := os.Getenv("CREATE_ADMIN_USER") == "0" || os.Getenv("CREATE_ADMIN_USER") == "false"
	webdavDataDir := os.Getenv("WEBDAV_DATA_DIR")
	return webdavPort, createNoAdminUser, webdavDataDir
}

func Read() error {
	initViper()
	path := getConfigurationPath()
	logging.Log.Info.Printf("Reading configuration from %s\n", path)
	v.SetConfigFile(filepath.Join(path, "config.yaml"))
	webdavPort, createNoAdminUser, webdavDataDir := ReadEnv()
	defaultCfg := GenerateDefaultConfig(webdavPort, createNoAdminUser, webdavDataDir)
	v.SetDefault("network", defaultCfg.Network)
	v.SetDefault("content", defaultCfg.Content)
	v.SetDefault("users", defaultCfg.Users)
	readErr := v.ReadInConfig()
	var configFileNotFoundError viper.ConfigFileNotFoundError
	if errors.As(readErr, &configFileNotFoundError) {
		logging.Log.Info.Println("No configuration file found, creating default configuration")
		v.WriteConfig()
	}
	return nil
}

func GenerateDefaultConfig(webdavPort string, createNoAdminUser bool, webdavDataDir string) Config {
	defaultConfig := DeepCopyConfig(configTemplate)
	if webdavPort != "" {
		defaultConfig.Network.Port = webdavPort
	}
	if createNoAdminUser {
		defaultConfig.Users = &map[string]User{}
	}
	if webdavDataDir != "" {
		defaultConfig.Content.Dir = webdavDataDir
	}
	return defaultConfig
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
