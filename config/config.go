package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var configTemplate *template.Template

// WriteConfigFile renders config using the template and writes it to configFilePath.
func WriteConfigFile(configFilePath string, config interface{}) error {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		return err
	}

	if err := os.WriteFile(configFilePath, buffer.Bytes(), 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func GetConfig(configFilePath string, defaultConfig interface{}) (interface{}, error) {
	filename := filepath.Base(configFilePath)
	fileType := filepath.Ext(filename)
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(filename)
	v.AddConfigPath(configFilePath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read in %s: %w", configFilePath, err)
	}

	if err := v.Unmarshal(&defaultConfig); err != nil {
		return nil, fmt.Errorf("error extracting app config: %w", err)
	}

	return &defaultConfig, nil
}
