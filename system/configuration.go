package system

//go:generate rice embed-go

import (
	"encoding/json"
	"os"
	"path"
)

//Configs contains application configurations for all application modes
type Configs struct {
	Debug   Config
	Release Config
	Test    Config
}

//Config contains application configuration for active application mode
type Config struct {
	Public        string `json:"public"`
	Domain        string `json:"domain"`
	SessionSecret string `json:"session_secret"`
	SignupEnabled bool   `json:"signup_enabled"` //always set to false in release mode (config.json)
	Database      DatabaseConfig
}

//DatabaseConfig contains database connection info
type DatabaseConfig struct {
	Host     string
	Name     string //database name
	User     string
	Password string
}

var (
	config *Config
)

//loadConfig unmarshals config for current application mode
func loadConfig(data []byte) {
	configs := &Configs{}
	err := json.Unmarshal(data, configs)
	if err != nil {
		panic(err)
	}
	switch GetMode() {
	case DebugMode:
		config = &configs.Debug
	case ReleaseMode:
		config = &configs.Release
	case TestMode:
		config = &configs.Test
	}
	if !path.IsAbs(config.Public) {
		workingDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		config.Public = path.Join(workingDir, config.Public)
	}
}

//GetConfig returns actual config
func GetConfig() *Config {
	return config
}

//PublicPath returns path to application public folder
func PublicPath() string {
	return config.Public
}

//UploadsPath returns path to public/uploads folder
func UploadsPath() string {
	return path.Join(config.Public, "uploads")
}
