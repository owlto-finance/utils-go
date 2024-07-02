package apollosdk

import (
	"fmt"
	"strings"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
)

type SDKConfig struct {
	AppID            string
	Cluster          string
	MetaAddr         string
	Namespaces       []string
	Secret           string
	IsBackupConfig   bool
	BackupConfigPath string
}

type ApolloSDK struct {
	client agollo.Client
	config SDKConfig
}

func NewApolloSDK(cfg SDKConfig) (*ApolloSDK, error) {
	clientConfig := &config.AppConfig{
		AppID:            cfg.AppID,
		Cluster:          cfg.Cluster,
		IP:               cfg.MetaAddr,
		IsBackupConfig:   cfg.IsBackupConfig,
		BackupConfigPath: cfg.BackupConfigPath,
		Secret:           cfg.Secret,
	}
	clientConfig.NamespaceName = strings.Join(cfg.Namespaces, config.Comma)

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return clientConfig, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start Apollo client: %w", err)
	}

	sdk := &ApolloSDK{
		client: client,
		config: cfg,
	}

	return sdk, nil
}

func (sdk *ApolloSDK) GetString(namespace, key string) (string, error) {
	config := sdk.client.GetConfig(namespace)
	if config == nil {
		return "", fmt.Errorf("namespace %s not found", namespace)
	}
	return config.GetValue(key), nil
}
