// factory.go
package provider

import (
	"fmt"

	"github.com/yourusername/smartdns/core"
	"github.com/yourusername/smartdns/config"
	"github.com/yourusername/smartdns/provider/local"
	"github.com/yourusername/smartdns/provider/tencent"
)

// New 根据配置创建SmartDNS实例
func New(providerType, configPath string) (core.SmartDNS, error) {
	switch providerType {
	case "local":
		cfg, err := config.Load[config.LocalConfig](configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load local config: %w", err)
		}
		return local.NewProvider(&local.Config{
			HostsPath:   cfg.HostsPath,
			StoragePath: cfg.StoragePath,
			Domain:      cfg.Domain,
		}), nil

	case "tencent":
		cfg, err := config.Load[config.TencentConfig](configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load tencent config: %w", err)
		}
		return tencent.NewProvider(&tencent.Config{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
			Region:    cfg.Region,
			Domain:    cfg.Domain,
		}), nil

	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}
