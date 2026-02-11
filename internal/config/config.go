package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env:"ADDR" env-required:"true"`
}

type Config struct {
	Env        string     `yaml:"env" env:"ENV" env-required:"true"`
	HTTPServer `yaml:"http_server"`

	DBHost     string `yaml:"db_host" env:"DB_HOST" env-required:"true"`
	DBPort     int    `yaml:"db_port" env:"DB_PORT" env-required:"true"`
	DBUser     string `yaml:"db_user" env:"DB_USER" env-required:"true"`
	DBPassword string `yaml:"db_password" env:"DB_PASSWORD" env-required:"true"`
	DBName     string `yaml:"db_name" env:"DB_NAME" env-required:"true"`
	DBSSLMode  string `yaml:"db_sslmode" env:"DB_SSLMODE" env-default:"disable"`
}

func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()
		if *flags != "" {
			configPath = *flags
		}
		if configPath == "" {
			log.Fatal("config path is required")
		}
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %s", err.Error())
	}
	return &cfg
}
