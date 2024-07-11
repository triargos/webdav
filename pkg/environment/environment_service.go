package environment

import "os"

type Service interface {
	Get(key string) string
	GetWithDefault(key string, defaultValue string) string
	GetBool(key string) bool
	Set(key string, value string) error
}

func NewOsEnvironmentService() Service {
	return &OsEnvironmentService{}
}

type OsEnvironmentService struct{}

func (s *OsEnvironmentService) Get(key string) string {
	return os.Getenv(key)
}

func (s *OsEnvironmentService) Set(key string, value string) error {
	return os.Setenv(key, value)
}

func (s *OsEnvironmentService) GetWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (s *OsEnvironmentService) GetBool(key string) bool {
	return os.Getenv(key) == "true" || os.Getenv(key) == "1"
}
