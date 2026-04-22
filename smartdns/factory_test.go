// factory_test.go
package smartdns

import (
	"os"
	"testing"
)

func TestNewLocalProvider(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := tmpDir + "/local.yaml"

	// 创建配置文件
	content := `
hosts_path: "/etc/hosts"
storage_path: "` + tmpDir + `/records.json"
domain: "smartdns.local"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	dns, err := New("local", configPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if dns == nil {
		t.Error("New() returned nil")
	}
}

func TestNewTencentProvider(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := tmpDir + "/tencent.yaml"

	content := `
secret_id: "test_id"
secret_key: "test_key"
region: "ap-guangzhou"
domain: "example.com"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	dns, err := New("tencent", configPath)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if dns == nil {
		t.Error("New() returned nil")
	}
}

func TestNewInvalidProvider(t *testing.T) {
	_, err := New("invalid", "/tmp/config.yaml")
	if err == nil {
		t.Error("Expected error for invalid provider type")
	}
}
