package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// SubDomain 从完整域名中提取子域名部分
func SubDomain(fullDomain, zoneDomain string) string {
	if fullDomain == zoneDomain {
		return "@"
	}
	return strings.TrimSuffix(fullDomain, "."+zoneDomain)
}

// FullDomain 从子域名构造完整域名
func FullDomain(sub, zoneDomain string) string {
	if sub == "@" {
		return zoneDomain
	}
	return sub + "." + zoneDomain
}

// GenerateCNAME 生成 CNAME，hash 拼接在子域名后
func GenerateCNAME(fullDomain, zoneDomain string) string {
	sub := SubDomain(fullDomain, zoneDomain)
	h := sha256.Sum256([]byte(zoneDomain + ":" + fullDomain))
	hash := hex.EncodeToString(h[:3])
	return sub + "-" + hash + "." + zoneDomain
}
