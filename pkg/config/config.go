package config

import (
	"github.com/triargos/webdav/pkg/fs"
	"github.com/triargos/webdav/pkg/logging"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Network *NetworkConfig   `yaml:"network"`
	Content *ContentConfig   `yaml:"content"`
	Users   *map[string]User `yaml:"users"`
}

type NetworkConfig struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
	Prefix  string `yaml:"prefix,omitempty"`
}

type ContentConfig struct {
	Dir string `yaml:"dir"`
}

type BasicAuthConfig struct {
	Realm string `yaml:"realm"`
}

type User struct {
	Password       string   `yaml:"password"`
	Root           string   `yaml:"root,omitempty"`
	SubDirectories []string `yaml:"sub_directories,omitempty"`
	Jail           bool     `yaml:"jail,omitempty"`
	Admin          bool     `yaml:"admin"`
}

func (c *Config) AddUser(username string, user User) {
	(*c.Users)[username] = user
}

func (u *User) UpdatePassword(password string) {
	u.Password = password
}

type LoggingConfig struct {
	Error  bool `yaml:"error"`
	Create bool `yaml:"create"`
	Read   bool `yaml:"read"`
	Update bool `yaml:"update"`
	Delete bool `yaml:"delete"`
}

type CORSConfig struct {
	Origin string `yaml:"origin"`
}

var Value *Config

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
var configDir = "./config/webdav"

func Init() error {
	isDocker := os.Getenv("DOCKER_ENABLED") == "1"
	if isDocker {
		configDir = "/etc/webdav/config"
	}
	config, err := ReadConfig()
	if err != nil {
		return err
	}
	Value = config
	return nil
}

func WriteDefaultConfig() error {
	return WriteConfig(&defaultConfig)
}

func ReadConfig() (*Config, error) {
	logging.Log.Info.Printf("Reading config from %s...", filepath.Join(configDir, "config.yaml"))
	if !fs.PathExists(filepath.Join(configDir, "config.yaml")) {
		//Write default config to file
		err := WriteConfig(&defaultConfig)
		if err != nil {
			return nil, err
		}
		return &defaultConfig, nil
	}
	file, err := os.Open(filepath.Join(configDir, "config.yaml"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func WriteConfig(config *Config) error {
	//Make the path if it doesn't exist
	if !fs.PathExists(configDir) {
		err := os.Mkdir(configDir, 0755)
		if err != nil {
			return err
		}
	}
	//Write config to file
	file, err := os.OpenFile(filepath.Join(configDir, "config.yaml"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	marshalled, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	_, err = file.Write(marshalled)
	if err != nil {
		return err
	}
	return nil
}
