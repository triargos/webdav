package config

import (
	"github.com/spf13/viper"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigurationPath(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		expectedPath string
	}{
		{"DockerEnabled", "1", "/etc/webdav"},
		{"DockerDisabled", "0", "./config"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DOCKER_ENABLED", tt.envValue)
			assert.Equal(t, tt.expectedPath, getConfigurationPath())
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name        string
		configSetup func()
		expected    *Config
	}{
		{
			name: "ValidConfig",
			configSetup: func() {
				v.Set("network::address", "127.0.0.1")
				v.Set("network::port", "9090")
				v.Set("content::dir", "/tmp/webdav")
				v.Set("users", map[string]User{
					"testuser": {
						Password:       "password",
						Admin:          false,
						Jail:           true,
						Root:           "/home/testuser",
						SubDirectories: []string{"subdir1", "subdir2"},
					},
				})
			},
			expected: &Config{
				Network: &NetworkConfig{
					Address: "127.0.0.1",
					Port:    "9090",
					Prefix:  "",
				},
				Content: &ContentConfig{
					Dir: "/tmp/webdav",
				},
				Users: &map[string]User{
					"testuser": {
						Password:       "password",
						Admin:          false,
						Jail:           true,
						Root:           "/home/testuser",
						SubDirectories: []string{"subdir1", "subdir2"},
					},
				},
			},
		},
		{
			name: "DefaultConfig",
			configSetup: func() {
				v = viper.NewWithOptions(viper.KeyDelimiter("::"))
				v.SetDefault("network", configTemplate.Network)
				v.SetDefault("content", configTemplate.Content)
				v.SetDefault("users", configTemplate.Users)
			},
			expected: &configTemplate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.configSetup()
			assert.Equal(t, tt.expected, Get())
		})
	}
}

func TestAddAndRemoveUser(t *testing.T) {
	tests := []struct {
		name             string
		initialUsers     map[string]User
		userToAdd        User
		usernameToAdd    string
		usernameToRemove string
		expectedUsers    map[string]User
	}{
		{
			name: "AddUser",
			initialUsers: map[string]User{
				"existinguser": {Password: "pass1", Admin: false},
			},
			userToAdd: User{
				Password: "newpass",
				Admin:    true,
			},
			usernameToAdd: "newuser",
			expectedUsers: map[string]User{
				"existinguser": {Password: "pass1", Admin: false},
				"newuser":      {Password: "newpass", Admin: true},
			},
		},
		{
			name: "RemoveUser",
			initialUsers: map[string]User{
				"existinguser": {Password: "pass1", Admin: false},
				"userToRemove": {Password: "pass2", Admin: true},
			},
			usernameToRemove: "userToRemove",
			expectedUsers: map[string]User{
				"existinguser": {Password: "pass1", Admin: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Users: &tt.initialUsers}
			Set(cfg)

			if tt.usernameToAdd != "" {
				AddUser(tt.usernameToAdd, tt.userToAdd)
			}

			if tt.usernameToRemove != "" {
				RemoveUser(tt.usernameToRemove)
			}

			result := Get()
			assert.Equal(t, tt.expectedUsers, *result.Users)
		})
	}
}

func TestWriteDefaultConfig(t *testing.T) {
	tests := []struct {
		name              string
		webdavPort        string
		createNoAdminUser bool
		webdavDataDir     string
		expectedNetwork   *NetworkConfig
		expectedContent   *ContentConfig
		expectedUsers     *map[string]User
	}{
		{
			name:              "All flags set",
			webdavPort:        "9090",
			createNoAdminUser: true,
			webdavDataDir:     "/tmp/webdav",
			expectedNetwork: &NetworkConfig{
				Address: "0.0.0.0",
				Port:    "9090",
				Prefix:  "/",
			},
			expectedContent: &ContentConfig{
				Dir: "/tmp/webdav",
			},
			expectedUsers: &map[string]User{},
		},
		{
			name:              "Create admin user not set",
			webdavPort:        "9090",
			createNoAdminUser: false,
			webdavDataDir:     "/tmp/webdav",
			expectedNetwork: &NetworkConfig{
				Address: "0.0.0.0",
				Port:    "9090",
				Prefix:  "/",
			},
			expectedContent: &ContentConfig{
				Dir: "/tmp/webdav",
			},
			expectedUsers: &map[string]User{
				"admin": {
					Password:       "admin",
					Admin:          true,
					Jail:           false,
					Root:           "/Users/admin",
					SubDirectories: []string{"documents"},
				},
			},
		},
		{
			name:              "Webdav data dir not set",
			webdavPort:        "9090",
			createNoAdminUser: false,
			webdavDataDir:     "",
			expectedNetwork: &NetworkConfig{
				Address: "0.0.0.0",
				Port:    "9090",
				Prefix:  "/",
			},
			expectedContent: &ContentConfig{
				Dir: "/var/webdav/data",
			},
			expectedUsers: &map[string]User{
				"admin": {
					Password:       "admin",
					Admin:          true,
					Jail:           false,
					Root:           "/Users/admin",
					SubDirectories: []string{"documents"},
				},
			},
		},
		{
			name:              "Port not set",
			webdavPort:        "",
			createNoAdminUser: false,
			webdavDataDir:     "/tmp/webdav",
			expectedNetwork: &NetworkConfig{
				Address: "0.0.0.0",
				Port:    "8080",
				Prefix:  "/",
			},
			expectedContent: &ContentConfig{
				Dir: "/tmp/webdav",
			},
			expectedUsers: &map[string]User{
				"admin": {
					Password:       "admin",
					Admin:          true,
					Jail:           false,
					Root:           "/Users/admin",
					SubDirectories: []string{"documents"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newConfig := GenerateDefaultConfig(tt.webdavPort, tt.createNoAdminUser, tt.webdavDataDir)
			assert.Equal(t, tt.expectedNetwork, newConfig.Network)
			assert.Equal(t, tt.expectedContent, newConfig.Content)
			assert.Equal(t, tt.expectedUsers, newConfig.Users)
		})
	}
}
