package config

type Config struct {
	Network  NetworkConfig   `mapstructure:"network"`
	Content  ContentConfig   `mapstructure:"content"`
	Users    map[string]User `mapstructure:"users"`
	Security SecurityConfig  `mapstructure:"security"`
}

type SecurityConfig struct {
	AuthType string `mapstructure:"authtype"`
}

type NetworkConfig struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
	Prefix  string `mapstructure:"prefix,omitempty"`
}

type ContentConfig struct {
	Dir            string   `mapstructure:"dir"`
	SubDirectories []string `mapstructure:"subdirectories,omitempty"`
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
		Dir:            "/var/webdav/data",
		SubDirectories: []string{"documents"},
	},
	Security: SecurityConfig{
		AuthType: "basic",
	},
	Users: map[string]User{},
}

func DeepCopyConfig(original Config) Config {
	newConfig := Config{
		Network: NetworkConfig{
			Address: original.Network.Address,
			Port:    original.Network.Port,
			Prefix:  original.Network.Prefix,
		},
		Security: SecurityConfig{
			AuthType: original.Security.AuthType,
		},
		Content: ContentConfig{
			Dir: original.Content.Dir,
		},
		Users: map[string]User{},
	}

	return newConfig
}

type EnvironmentConfig struct {
	WebdavPort    string
	WebdavDataDir string
	AuthType      string
}
