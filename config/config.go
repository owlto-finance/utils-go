package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/owlto-finance/utils-go/system"
	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

// WriteConfigFile renders config using the template and writes it to configFilePath.
func WriteConfigFile(configFilePath string, config interface{}) error {
	fileType := strings.TrimPrefix(filepath.Ext(configFilePath), ".")

	configMap := make(map[string]interface{})
	if err := mapstructure.Decode(config, &configMap); err != nil {
		return fmt.Errorf("error decoding struct to map: %w", err)

	}
	err := system.MakedirAll(filepath.Dir(configFilePath))
	if err != nil {
		return err
	}

	var data []byte
	if fileType == "toml" {
		// Marshal the map to TOML format
		data, err = toml.Marshal(configMap)
		if err != nil {
			return fmt.Errorf("error marshaling map to TOML: %w", err)
		}
	} else {
		return fmt.Errorf("unsupport config file type: %s, only toml is support", fileType)
	}

	if err := os.WriteFile(configFilePath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func GetConfig(configFilePath string, defaultConfig interface{}) error {
	filename := filepath.Base(configFilePath)
	fileType := strings.TrimPrefix(filepath.Ext(configFilePath), ".")
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(filename)
	v.AddConfigPath(filepath.Dir(configFilePath))

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read in %s: %w", configFilePath, err)
	}

	if err := v.Unmarshal(&defaultConfig); err != nil {
		return fmt.Errorf("error extracting app config: %w", err)
	}

	return nil
}
