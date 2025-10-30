package config

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configFile = "tms.config.json"
)

type Config struct {
	Url                               string `json:"url" env-required:"true" env:"TMS_URL"`
	Token                             string `json:"privateToken" env-required:"true" env:"TMS_PRIVATE_TOKEN"`
	ProjectId                         string `json:"projectId" env-required:"true" env:"TMS_PROJECT_ID"`
	ConfigurationId                   string `json:"configurationId" env-required:"true" env:"TMS_CONFIGURATION_ID"`
	TestRunId                         string `json:"testRunId" env:"TMS_TEST_RUN_ID"`
	TestRunName                       string `json:"testRunName" env:"TMS_TEST_RUN_NAME"`
	AdapterMode                       string `json:"adapterMode" env:"TMS_ADAPTER_MODE" env-default:"0"`
	IsDebug                           bool   `json:"isDebug" env:"TMS_IS_DEBUG" env-default:"false"`
	AutomaticCreationTestCases        bool   `json:"automaticCreationTestCases" env:"TMS_AUTOMATIC_CREATION_TEST_CASES" env-default:"false"`
	AutomaticUpdationLinksToTestCases bool   `json:"automaticUpdationLinksToTestCases" env:"TMS_AUTOMATIC_UPDATION_LINKS_TO_TEST_CASES" env-default:"false"`
	CertValidation                    bool   `json:"certValidation" env:"TMS_CERT_VALIDATION" env-default:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("TMS_CONFIG_FILE")
	if configPath == "" {
		// If not provided, try to find it in current working directory
		configPath = getDefaultConfigPath()
		// Check if file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("TMS_CONFIG_FILE is not set and config file not found in the current directory: %s", configFile)
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	if err := cleanenv.UpdateEnv(&cfg); err != nil {
		log.Fatalf("cannot update config: %s", err)
	}

	validateConfig(cfg)

	return &cfg
}

func validateConfig(cfg Config) {
	_, err := url.ParseRequestURI(cfg.Url)
	if err != nil {
		panic("Url is invalid")
	}

	if cfg.Token == "null" || cfg.Token == "" {
		panic("Private token is invalid")
	}

	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")

	if !r.MatchString(cfg.ProjectId) {
		panic("Project id is invalid")
	}

	if !r.MatchString(cfg.ConfigurationId) {
		panic("Configuration id is invalid")
	}

	adapterMode, err := strconv.Atoi(cfg.AdapterMode)
	if err != nil {
		panic("Adapter mode is invalid")
	}

	// adapter mode 2 - test run id is created automatically
	if adapterMode == 2 {
		// do nothing
	} else if adapterMode == 1 {
		if !r.MatchString(cfg.TestRunId) {
			panic("Adapter works in mode 1. Config should contains valid test run id")
		}
	} else if adapterMode == 0 {
		if !r.MatchString(cfg.TestRunId) {
			panic("Adapter works in mode 0. Config should contains valid test run id")
		}
	} else {
		panic("Adapter mode is invalid")
	}
}

func getDefaultConfigPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get current working directory: %s", err)
	}

	return filepath.Join(cwd, configFile)
}
