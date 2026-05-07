// provider/tencent/provider_test.go
package tencent_test

import (
	"strings"
	"testing"

	"github.com/wangbo2295/gadns/core"
	"github.com/wangbo2295/gadns/provider/tencent"
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

	var _ core.CNAMEProvider = provider
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

	cname := provider.GenerateCNAME("app.example.com")
	// 格式：app-<6位hex哈希>.example.com
	if !strings.HasPrefix(cname, "app-") || !strings.HasSuffix(cname, ".example.com") {
		t.Errorf("GenerateCNAME() = %v, want app-<hash>.example.com", cname)
	}

	// 相同输入应产生相同 CNAME（确定性）
	cname2 := provider.GenerateCNAME("app.example.com")
	if cname != cname2 {
		t.Errorf("GenerateCNAME() not deterministic: %v != %v", cname, cname2)
	}

	// 不同 name 应产生不同 CNAME
	cname3 := provider.GenerateCNAME("web.example.com")
	if cname == cname3 {
		t.Errorf("GenerateCNAME() should differ for different names")
	}
}
