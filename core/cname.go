// core/cname.go - CNAMEProvider 核心接口和类型定义
package core

// CNAMEProvider 提供 CNAME 生成和 DNS 记录管理能力
type CNAMEProvider interface {
	Create(name string, ips []string) (*Record, error)
	Update(name string, ips []string) (*Record, error)
	Get(name string) (*Record, error)
	List() ([]*Record, error)
	Delete(name string) error
}

// Record DNS 记录信息
type Record struct {
	Name  string   // 完整域名
	CNAME string   // 生成的 CNAME
	IPs   []string // 原始 IP 列表
}
