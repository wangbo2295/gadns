// factory.go
package provider

import (
	"fmt"

	"github.com/wangbo2295/gadns/config"
	"github.com/wangbo2295/gadns/core"
	"github.com/wangbo2295/gadns/provider/tencent"
)

// New 根据配置创建 CNAMEProvider 实例（目前仅支持腾讯云）
func New(providerType, configPath string) (core.CNAMEProvider, error) {
	switch providerType {
	case "tencent":
		cfg, err := config.Load[config.TencentConfig](configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load tencent config: %w", err)
		}
		provider, err := tencent.NewProvider(&tencent.Config{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
			Region:    cfg.Region,
			Domain:    cfg.Domain,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create tencent provider: %w", err)
		}
		return provider, nil

	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}
