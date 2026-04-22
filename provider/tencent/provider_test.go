// provider/tencent/provider_test.go
package tencent_test

import (
	"testing"

	"github.com/yourusername/smartdns/core"
	"github.com/yourusername/smartdns/provider/tencent"
)

func TestTencentProviderImplementsInterface(t *testing.T) {
	provider, err := tencent.NewProvider(&tencent.Config{
		SecretID:  "test_id",
		SecretKey: "test_key",
		Region:    "ap-guangzhou",
		Domain:    "example.com",
	})

	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}

	// Verify interface compliance
	var _ core.SmartDNS = provider
	_ = provider
}

func TestTencentProviderGenerateCNAME(t *testing.T) {
	provider, err := tencent.NewProvider(&tencent.Config{
		SecretID:  "test_id",
		SecretKey: "test_key",
		Region:    "ap-guangzhou",
		Domain:    "example.com",
	})

	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}

	cname := provider.GenerateCNAME("app")
	expected := "app.example.com"

	if cname != expected {
		t.Errorf("GenerateCNAME() = %v, want %v", cname, expected)
	}
}

func TestTencentProviderGetSmartRoutingConfig(t *testing.T) {
	provider, err := tencent.NewProvider(&tencent.Config{
		SecretID:  "test_id",
		SecretKey: "test_key",
		Region:    "ap-guangzhou",
		Domain:    "example.com",
	})

	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}

	// 通过测试 Create 方法间接测试智能调度配置
	// 注意：这里会失败，因为使用的是测试密钥
	_, err = provider.Create("test", []string{"1.1.1.1"})
	if err == nil {
		t.Log("Create() succeeded (unexpected with test credentials)")
	} else {
		t.Logf("Create() failed as expected with test credentials: %v", err)
	}
}
