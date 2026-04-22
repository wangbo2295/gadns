// config/loader.go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TencentConfig 腾讯云配置
type TencentConfig struct {
	SecretID  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
	Region    string `yaml:"region"`
	Domain    string `yaml:"domain"`
}

// LocalConfig 本地实现配置
type LocalConfig struct {
	HostsPath   string `yaml:"hosts_path"`
	StoragePath string `yaml:"storage_path"`
	Domain      string `yaml:"domain"`
}

// Load 从文件加载配置
func Load[T any](path string) (*T, error) {
	// 展开波浪线路径
	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[1:])
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg T
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}
