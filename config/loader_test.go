// config/loader_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTencentConfig(t *testing.T) {
	// 创建临时配置文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "tencent.yaml")

	content := `
secret_id: "test_secret_id"
secret_key: "test_secret_key"
region: "ap-guangzhou"
domain: "example.com"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load[TencentConfig](configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.SecretID != "test_secret_id" {
		t.Errorf("SecretID = %v, want test_secret_id", cfg.SecretID)
	}
	if cfg.SecretKey != "test_secret_key" {
		t.Errorf("SecretKey = %v, want test_secret_key", cfg.SecretKey)
	}
	if cfg.Region != "ap-guangzhou" {
		t.Errorf("Region = %v, want ap-guangzhou", cfg.Region)
	}
	if cfg.Domain != "example.com" {
		t.Errorf("Domain = %v, want example.com", cfg.Domain)
	}
}

func TestLoadInvalidPath(t *testing.T) {
	_, err := Load[TencentConfig]("/nonexistent/config.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}
