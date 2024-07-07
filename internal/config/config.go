package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env             string        `yaml:"env" env:"ENV" env-default:"local"`
	HTTP            HTTP          `yaml:"http"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
	/*GRPC     GRPC     `yaml:"grpc"`
	Rabbitmq Rabbitmq `yaml:"rabbitmq"`
	MongoDB  MongoDB  `yaml:"mongodb"`*/
}

type HTTP struct {
	Address     string        `yaml:"address" env:"HTTP_ADDRESS" env-default:"localhost:5000"`
	Domain      string        `yaml:"domain" env:"HTTP_DOMAIN"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"4s"`
}

type GRPC struct {
	Port    int           `yaml:"port" env:"GRPC_PORT"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT"`
}

type MongoDB struct {
	URI          string        `yaml:"uri" env:"MONGODB_URI"`
	PingTimeout  time.Duration `yaml:"ping_timeout" env:"MONGODB_PING_TIMEOUT" env-default:"10s"`
	DatabaseName string        `yaml:"database_name" env:"MONGODB_DATABASE_NAME" env-default:"uniposts"`
}

type Rabbitmq struct {
	User     string `yaml:"user" env:"RABBITMQ_USER"`
	Password string `yaml:"password" env:"RABBITMQ_PASSWORD"`
	Host     string `yaml:"host" env:"RABBITMQ_HOST"`
	Port     string `yaml:"port" env:"RABBITMQ_PORT"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		return mustLoadFromEnv()
	}

	return mustLoadByPath(path)
}

func mustLoadByPath(configPath string) *Config {
	cfg, err := loadByPath(configPath)
	if err != nil {
		panic(err)
	}

	return cfg
}

func loadByPath(configPath string) (*Config, error) {
	var cfg Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("there is no config file: %w", err)
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}

func mustLoadFromEnv() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Env empty")
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string

	if flag.Lookup("config") == nil {
		flag.StringVar(&res, "config", "", "path to config file")
		flag.Parse()
	}

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
