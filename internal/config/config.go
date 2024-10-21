package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func (cfg *Config) IsLocal() bool {
	return cfg.Env == "local"
}

func (cfg *Config) IsProd() bool {
	return cfg.Env == "prod"
}

func MustLoad() *Config {
	path := fetchConfig()
	if path == "" {
		panic("config path empty")
	}
	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file doesn't exist " + path)
	}

	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("couldn't load config")
	}
	return &cfg
}

func fetchConfig() string {
	var ret string

	flag.StringVar(&ret, "config", "", "config file")
	flag.Parse()

	if ret == "" {
		ret = os.Getenv("CONFIG_PATH")
	}
	return ret
}
