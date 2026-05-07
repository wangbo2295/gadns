// provider/tencent/client_test.go
package tencent

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(&Config{
		SecretID:  "test_id",
		SecretKey: "test_key",
		Region:    "ap-guangzhou",
		Domain:    "example.com",
	})

	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Error("NewClient() returned nil")
	}
	if client.domain != "example.com" {
		t.Errorf("domain = %v, want example.com", client.domain)
	}
}
