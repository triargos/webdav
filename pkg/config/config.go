package config

import (
	"github.com/triargos/webdav/pkg/fs"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Network NetworkConfig   `yaml:"network"`
	TLS     TLSConfig       `yaml:"tls"`
	Content ContentConfig   `yaml:"content"`
	Users   map[string]User `yaml:"users"`
}

type NetworkConfig struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
	Prefix  string `yaml:"prefix,omitempty"`
}

type TLSConfig struct {
	KeyFile  string `yaml:"keyFile"`
	CertFile string `yaml:"certFile"`
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

var Value Config

var defaultConfig = Config{
	Network: NetworkConfig{
		Address: "0.0.0.0",
		Port:    "8080",
		Prefix:  "/",
	},
	TLS: TLSConfig{
		KeyFile:  "server.key",
		CertFile: "server.crt",
	},
	Content: ContentConfig{
		Dir: "/var/webdav/data",
	},
	Users: map[string]User{
		"admin": {
			Password: "admin",
			Jail:     false,
			Root:     "/Users/admin",
		},
	},
}
var configPath = "/etc/webdav/config.yaml"

func Init() error {
	config, err := ReadConfig()
	if err != nil {
		return err
	}
	Value = *config
	return nil
}

func ReadConfig() (*Config, error) {
	if !fs.PathExists(configPath) {
		//Write default config to file
		err := WriteConfig(&defaultConfig)
		if err != nil {
			return nil, err
		}
		return &defaultConfig, nil
	}
	file, err := os.Open(configPath)
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
	if !fs.PathExists("/etc/webdav") {
		err := os.Mkdir("/etc/webdav", 0755)
		if err != nil {
			return err
		}
	}
	//Write config to file
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0644)
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
