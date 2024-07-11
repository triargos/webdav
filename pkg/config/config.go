package config

type Config struct {
	Network NetworkConfig   `mapstructure:"network"`
	Content ContentConfig   `mapstructure:"content"`
	Users   map[string]User `mapstructure:"users"`
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

var configTemplate = Config{
	Network: NetworkConfig{
		Address: "0.0.0.0",
		Port:    "8080",
		Prefix:  "/",
	},
	Content: ContentConfig{
		Dir: "/var/webdav/data",
	},
	Users: map[string]User{
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
		Network: NetworkConfig{
			Address: original.Network.Address,
			Port:    original.Network.Port,
			Prefix:  original.Network.Prefix,
		},
		Content: ContentConfig{
			Dir: original.Content.Dir,
		},
		Users: map[string]User{},
	}

	for k, v := range original.Users {
		newUser := v
		(newConfig.Users)[k] = newUser
	}

	return newConfig
}

type EnvironmentConfig struct {
	WebdavPort        string
	CreateNoAdminUser bool
	WebdavDataDir     string
}
