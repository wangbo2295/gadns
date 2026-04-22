// provider/tencent/provider_test.go
package tencent_test

import (
	"testing"

	"github.com/yourusername/smartdns/smartdns/provider/tencent"
	"github.com/yourusername/smartdns/types"
)

func TestTencentProviderImplementsInterface(t *testing.T) {
	provider := tencent.NewProvider(&tencent.Config{
		SecretID:  "test_id",
		SecretKey: "test_key",
		Region:    "ap-guangzhou",
		Domain:    "example.com",
	})

	// 验证接口实现
	var _ types.SmartDNS = provider
}

func TestTencentProviderGenerateCNAME(t *testing.T) {
	provider := tencent.NewProvider(&tencent.Config{
		SecretID:  "test_id",
		SecretKey: "test_key",
		Region:    "ap-guangzhou",
		Domain:    "example.com",
	})

	cname := provider.GenerateCNAME("app")
	expected := "app.example.com"

	if cname != expected {
		t.Errorf("GenerateCNAME() = %v, want %v", cname, expected)
	}
}
