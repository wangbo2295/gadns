// provider/tencent/provider_test.go
package tencent_test

import (
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
