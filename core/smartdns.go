// core/smartdns.go - SmartDNS 核心接口和类型定义
package core

// SmartDNS 组件接口
type SmartDNS interface {
	Create(name string, ips []string) (*Record, error)
	Update(name string, ips []string) (*Record, error)
	Get(name string) (*Record, error)
	List() ([]*Record, error)
	Delete(name string) error
}

// Record DNS记录信息
type Record struct {
	Name  string   // 记录名称，用于后续Update/Delete
	CNAME string   // 生成的CNAME，调用者用这个做DNS解析
	IPs   []string // 原始IP列表（支持列表、范围、段格式）
}
