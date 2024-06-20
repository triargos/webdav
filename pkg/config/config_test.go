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
				v.SetDefault("network", defaultConfig.Network)
				v.SetDefault("content", defaultConfig.Content)
				v.SetDefault("users", defaultConfig.Users)
			},
			expected: &defaultConfig,
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
