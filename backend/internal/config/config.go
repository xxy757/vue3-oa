package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Redis    RedisConfig    `yaml:"redis"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type RedisConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("OA_CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	localPath := "configs/config.local.yaml"
	if localData, err := os.ReadFile(localPath); err == nil {
		var localCfg Config
		if err := yaml.Unmarshal(localData, &localCfg); err == nil {
			cfg.Merge(&localCfg)
		}
	}

	if v := os.Getenv("OA_DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}

	return &cfg, nil
}

func (c *Config) Merge(other *Config) {
	if other.Server.Port != 0 {
		c.Server.Port = other.Server.Port
	}
	if other.Server.Mode != "" {
		c.Server.Mode = other.Server.Mode
	}
	if other.Database.Host != "" {
		c.Database.Host = other.Database.Host
	}
	if other.Database.Port != 0 {
		c.Database.Port = other.Database.Port
	}
	if other.Database.User != "" {
		c.Database.User = other.Database.User
	}
	if other.Database.Password != "" {
		c.Database.Password = other.Database.Password
	}
	if other.Database.DBName != "" {
		c.Database.DBName = other.Database.DBName
	}
	if other.Database.MaxIdleConns != 0 {
		c.Database.MaxIdleConns = other.Database.MaxIdleConns
	}
	if other.Database.MaxOpenConns != 0 {
		c.Database.MaxOpenConns = other.Database.MaxOpenConns
	}
	if other.JWT.Secret != "" {
		c.JWT.Secret = other.JWT.Secret
	}
	if other.JWT.ExpireHours != 0 {
		c.JWT.ExpireHours = other.JWT.ExpireHours
	}
}
