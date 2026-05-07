package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// GenerateCNAME 生成 CNAME：将域名中的 . 替换为 -，拼接 hash，以 zoneDomain 结尾
// 例如 app.doerhh.cn → app-doerhh-cn-a1b2c3.doerhh.cn
func GenerateCNAME(fullDomain, zoneDomain string) string {
	safe := strings.ReplaceAll(fullDomain, ".", "-")
	h := sha256.Sum256([]byte(zoneDomain + ":" + fullDomain))
	hash := hex.EncodeToString(h[:3])
	return safe + "-" + hash + "." + zoneDomain
}
