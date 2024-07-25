package config

type Config struct {
	Network  NetworkConfig   `yaml:"network"`
	Content  ContentConfig   `yaml:"content"`
	Users    map[string]User `yaml:"users"`
	Security SecurityConfig  `yaml:"security"`
}

type SecurityConfig struct {
	AuthType string `yaml:"authtype"`
}

type NetworkConfig struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
	Prefix  string `yaml:"prefix,omitempty"`
}

type ContentConfig struct {
	Dir            string   `yaml:"dir"`
	SubDirectories []string `yaml:"subdirectories,omitempty"`
}

type User struct {
	Password       string   `yaml:"password"`
	Root           string   `yaml:"root,omitempty"`
	SubDirectories []string `yaml:"subdirectories,omitempty"`
	Jail           bool     `yaml:"jail,omitempty"`
	Admin          bool     `yaml:"admin"`
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
