// provider/tencent/config.go
package tencent

// Config 腾讯云Provider配置
type Config struct {
	SecretID  string // 腾讯云 API 密钥 ID
	SecretKey string // 腾讯云 API 密钥 Key
	Region    string // 地域（默认：ap-guangzhou）
	Domain    string // DNSPod 域名
}

// DefaultRegion 默认地域
const DefaultRegion = "ap-guangzhou"

