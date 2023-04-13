package dnsprovider

import (
	"context"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type DnsProvider interface {
	UpdateRecord(ctx context.Context, record, zone, ip string, TTL int64) error
}

func ConfigureDnsProvider(ctx context.Context, provider, configFile string) (DnsProvider, error) {
	if provider == "aws" {
		config, err := ReadConfig(configFile, provider)
		if err != nil {
			return nil, err
		}

		accessKeyId, foundAccessKey := config[AWS_ACCESS_KEY]
		if !foundAccessKey {
			return nil, fmt.Errorf("AWS Provider missing configuration for %q", AWS_ACCESS_KEY)
		}

		secretAccessKey, foundSecretKey := config[AWS_SECRET_KEY]
		if !foundSecretKey {
			return nil, fmt.Errorf("AWS provider missing configuration for %q", AWS_SECRET_KEY)
		}

		aws, err := NewAws(ctx, accessKeyId, secretAccessKey)
		return aws, err
	}

	return nil, fmt.Errorf("provider %q not supported", provider)
}

func ReadConfig(configFile, provider string) (map[string]string, error) {
	var err error
	var configData []byte
	var providerConfig map[string]string
	var providerFound bool

	// Read config file
	configData, err = ioutil.ReadFile(configFile)
	if err != nil {
		return providerConfig, err
	}

	config := make(map[string]map[string]string)
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		return providerConfig, err
	}

	providerConfig, providerFound = config[provider]
	if !providerFound {
		return providerConfig, fmt.Errorf("provider %q not found in config %q", provider, configFile)
	}

	return providerConfig, err
}
