// cmd/root.go - 根命令和公共逻辑
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wangbo2295/gadns/config"
	"github.com/wangbo2295/gadns/core"
	"github.com/wangbo2295/gadns/provider"
)

var (
	cfgFile      string
	providerType string

	// Version 信息（由 ldflags 设置）
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "gadns",
	Short: "GADNS - DNS 记录管理工具",
	Long: `GADNS 根据 IP 集合生成 CNAME 记录，通过腾讯云 DNSPod 管理 DNS 记录。

每个 IP 创建一条 A 记录，多 IP 时自动分配权重实现负载均衡。`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, GitCommit, BuildTime),
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件路径 (默认 ~/.gadns/tencent.yaml)")
	rootCmd.PersistentFlags().StringVarP(&providerType, "provider", "p", "tencent", "Provider 类型 (tencent / noop)")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// Execute 执行命令
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func getConfigPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	return "~/.gadns/tencent.yaml"
}

func newProvider() (core.CNAMEProvider, error) {
	return provider.New(providerType, getConfigPath())
}

func resolveDomain() string {
	cfg, err := config.Load[config.TencentConfig](getConfigPath())
	if err != nil {
		return ""
	}
	return cfg.Domain
}

func fullName(name string) string {
	domain := resolveDomain()
	if domain == "" || strings.Contains(name, ".") {
		return name
	}
	return name + "." + domain
}
