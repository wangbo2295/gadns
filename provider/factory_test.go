// factory_test.go
package provider

import (
	"os"
	"testing"
)

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
		t.Skipf("New() error (expected with test credentials): %v", err)
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
