package config

import (
	"fmt"
	"github.com/triargos/webdav/pkg/environment"
	"github.com/triargos/webdav/pkg/fs"
	"github.com/triargos/webdav/pkg/helper"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
)

type Service interface {
	Get() *Config
	Set(cfg *Config)
	Write() error
	Read() *Config
	Reset() error

	CreateConfigDirectory() error

	GenerateDefault(config EnvironmentConfig) Config
	AddUser(username string, user User)
	UpdateUser(username string, user User)
	RemoveUser(username string)

	readEnvironmentConfig() EnvironmentConfig
}

type ConfigService struct {
	environmentService environment.Service
	fileSystemHandler  fs.Service
}

func NewConfigService(environmentService environment.Service, fileSystemHandler fs.Service) Service {
	service := ConfigService{environmentService: environmentService, fileSystemHandler: fileSystemHandler}
	createDirectoryErr := service.CreateConfigDirectory()
	if createDirectoryErr != nil {
		slog.Error("Failed to create config directory", "error", createDirectoryErr)
		os.Exit(1)
	}
	service.Read()
	return &service
}

var currentConfig *Config

func (s *ConfigService) Get() *Config {
	if currentConfig != nil {
		return currentConfig
	}
	return s.Read()
}

func (s *ConfigService) Set(cfg *Config) {
	currentConfig = cfg
}

func (s *ConfigService) Write() error {
	configPath := s.getConfigurationPath()
	marshalled, marshalError := yaml.Marshal(s.Get())
	if marshalError != nil {
		return fmt.Errorf("failed to marshal config: %w", marshalError)
	}
	writeFileErr := s.fileSystemHandler.WriteFileContent(configPath, marshalled, 0644)
	if writeFileErr != nil {
		return fmt.Errorf("failed to write config file: %w", writeFileErr)
	}
	return nil
}

func (s *ConfigService) Read() *Config {
	configPath := s.getConfigurationPath()
	environmentConfig := s.readEnvironmentConfig()

	defaultConfig := s.GenerateDefault(environmentConfig)
	currentConfig = &defaultConfig
	fileContents, readFileErr := s.fileSystemHandler.ReadFileContent(configPath)
	unmarshalErr := yaml.Unmarshal(fileContents, currentConfig)

	if readFileErr != nil || unmarshalErr != nil {
		slog.Info("Config file not found, writing default config")
		writeDefaultConfigErr := s.Write()
		if writeDefaultConfigErr != nil {
			slog.Error("Failed to write default config", "error", writeDefaultConfigErr)
			os.Exit(1)
		}
	}
	return s.Get()
}

func (s *ConfigService) GenerateDefault(environmentConfig EnvironmentConfig) Config {
	defaultConfig := DeepCopyConfig(configTemplate)
	if environmentConfig.WebdavPort != "" {
		defaultConfig.Network.Port = environmentConfig.WebdavPort
	}
	if environmentConfig.WebdavDataDir != "" {
		defaultConfig.Content.Dir = environmentConfig.WebdavDataDir
	}
	if environmentConfig.AuthType != "" && helper.ValidateAuthType(environmentConfig.AuthType) {
		defaultConfig.Security.AuthType = environmentConfig.AuthType
	}
	return defaultConfig
}

func (s *ConfigService) getConfigurationPath() string {
	defaultConfigPath := "/etc/webdav/"
	if s.environmentService.GetBool("DOCKER_ENABLED") == false {
		defaultConfigPath = "./config"
	}
	return filepath.Join(defaultConfigPath, "config.yaml")
}

func (s *ConfigService) AddUser(username string, user User) {
	currentUsers := s.Get().Users
	currentUsers[username] = user
	currentConfig.Users = currentUsers
}

func (s *ConfigService) RemoveUser(username string) {
	currentUsers := s.Get().Users
	delete(currentUsers, username)
	currentConfig.Users = currentUsers
}

func (s *ConfigService) UpdateUser(username string, user User) {
	s.AddUser(username, user)
}

func (s *ConfigService) readEnvironmentConfig() EnvironmentConfig {
	webdavPort := s.environmentService.Get("WEBDAV_PORT")
	webdavDataDir := s.environmentService.Get("WEBDAV_DATA_DIR")
	authType := s.environmentService.Get("AUTH_TYPE")
	return EnvironmentConfig{
		WebdavPort:    webdavPort,
		WebdavDataDir: webdavDataDir,
		AuthType:      authType,
	}
}

func (s *ConfigService) CreateConfigDirectory() error {
	configPath := s.getConfigurationPath()
	return os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
}

func (s *ConfigService) Reset() error {
	environmentConfig := s.readEnvironmentConfig()
	defaultConfig := s.GenerateDefault(environmentConfig)
	s.Set(&defaultConfig)
	return s.Write()
}
