package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configFile = "tms.config.json"
)

type Config struct {
	Url             string `json:"url" env-required:"true" env:"TMS_URL"`
	Token           string `json:"privateToken" env-required:"true" env:"TMS_PRIVATE_TOKEN"`
	ProjectId       string `json:"projectId" env-required:"true" env:"TMS_PROJECT_ID"`
	ConfigurationId string `json:"configurationId" env-required:"true" env:"TMS_CONFIGURATION_ID"`
	TestRunId       string `json:"testRunId" env-required:"true" env:"TMS_TEST_RUN_ID"`
	IsDebug         bool   `json:"isDebug" env:"TMS_IS_DEBUG" env-default:"false"`
}

func MustLoad() *Config {
	configPath := os.Getenv("TMS_CONFIG_FILE")
	if configPath == "" {
		log.Fatal("TMS_CONFIG_FILE is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
