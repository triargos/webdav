package config

import (
	"github.com/spf13/viper"
	"github.com/triargos/webdav/pkg/environment"
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

type ViperConfigService struct {
	viperInstance      *viper.Viper
	environmentService environment.Service
}

func NewViperConfigService(environmentService environment.Service) Service {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigType("yaml")
	service := ViperConfigService{viperInstance: v, environmentService: environmentService}
	createDirectoryErr := service.CreateConfigDirectory()
	if createDirectoryErr != nil {
		slog.Error("Failed to create config directory", "error", createDirectoryErr)
		os.Exit(1)
	}
	service.Read()
	return &service
}

func (s *ViperConfigService) Get() *Config {
	cfg := &Config{}
	if err := s.viperInstance.Unmarshal(&cfg); err != nil {
		envConfig := s.readEnvironmentConfig()
		defaultConfig := s.GenerateDefault(envConfig)
		return &defaultConfig
	}
	return cfg
}

func (s *ViperConfigService) Set(cfg *Config) {
	s.viperInstance.Set("network", cfg.Network)
	s.viperInstance.Set("content", cfg.Content)
	s.viperInstance.Set("users", cfg.Users)
}

func (s *ViperConfigService) Write() error {
	return s.viperInstance.WriteConfig()
}

func (s *ViperConfigService) Read() *Config {
	configPath := s.getConfigurationPath()
	s.viperInstance.SetConfigFile(filepath.Join(configPath, "config.yaml"))
	environmentConfig := s.readEnvironmentConfig()

	defaultConfig := s.GenerateDefault(environmentConfig)
	s.viperInstance.SetDefault("network", defaultConfig.Network)
	s.viperInstance.SetDefault("content", defaultConfig.Content)
	s.viperInstance.SetDefault("users", defaultConfig.Users)
	readFileErr := s.viperInstance.ReadInConfig()
	if readFileErr != nil {
		slog.Info("Config file not found, writing default config")
		writeDefaultConfigErr := s.viperInstance.WriteConfig()
		if writeDefaultConfigErr != nil {
			slog.Error("Failed to write default config", "error", writeDefaultConfigErr)
			os.Exit(1)
		}
	}
	return s.Get()
}

func (s *ViperConfigService) GenerateDefault(environmentConfig EnvironmentConfig) Config {
	defaultConfig := DeepCopyConfig(configTemplate)
	if environmentConfig.WebdavPort != "" {
		defaultConfig.Network.Port = environmentConfig.WebdavPort
	}
	if environmentConfig.CreateNoAdminUser {
		defaultConfig.Users = map[string]User{}
	}
	if environmentConfig.WebdavDataDir != "" {
		defaultConfig.Content.Dir = environmentConfig.WebdavDataDir
	}
	return defaultConfig
}

func (s *ViperConfigService) getConfigurationPath() string {
	defaultConfigPath := "/etc/webdav/"
	if s.environmentService.GetBool("DOCKER_ENABLED") == false {
		defaultConfigPath = "./config"
	}
	return defaultConfigPath
}

func (s *ViperConfigService) AddUser(username string, user User) {
	currentUsers := s.Get().Users
	currentUsers[username] = user
	s.viperInstance.Set("users", currentUsers)
}

func (s *ViperConfigService) RemoveUser(username string) {
	currentUsers := s.Get().Users
	delete(currentUsers, username)
	s.viperInstance.Set("users", currentUsers)
}

func (s *ViperConfigService) UpdateUser(username string, user User) {
	s.AddUser(username, user)
}

func (s *ViperConfigService) readEnvironmentConfig() EnvironmentConfig {
	webdavPort := s.environmentService.Get("WEBDAV_PORT")
	noAdminUser := s.environmentService.GetBool("CREATE_ADMIN_USER")
	webdavDataDir := s.environmentService.Get("WEBDAV_DATA_DIR")
	return EnvironmentConfig{
		WebdavPort:        webdavPort,
		CreateNoAdminUser: noAdminUser,
		WebdavDataDir:     webdavDataDir,
	}
}

func (s *ViperConfigService) CreateConfigDirectory() error {
	configPath := s.getConfigurationPath()
	return os.MkdirAll(configPath, os.ModePerm)
}

func (s *ViperConfigService) Reset() error {
	environmentConfig := s.readEnvironmentConfig()
	defaultConfig := s.GenerateDefault(environmentConfig)
	s.Set(&defaultConfig)
	return s.Write()
}
