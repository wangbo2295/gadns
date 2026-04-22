// provider/tencent/config.go
package tencent

// Config 腾讯云Provider配置
type Config struct {
	SecretID  string
	SecretKey string
	Region    string
	Domain    string
}
