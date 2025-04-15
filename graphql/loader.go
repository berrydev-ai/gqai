package graphql

import (
	"fmt"
	_ "os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Schema    string `yaml:"schema"`
	Documents string `yaml:"documents"`

	// The absolute path to the config directory
	BaseDir string
}

func LoadConfigAt(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	baseDir := filepath.Dir(path)
	config.BaseDir = baseDir

	// Convert relative doc path to absolute
	if !filepath.IsAbs(config.Documents) {
		config.Documents = filepath.Join(baseDir, config.Documents)
	}

	return &config, nil
}
