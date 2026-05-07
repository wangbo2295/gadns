// provider/noop/provider.go - 内网/测试用 Noop Provider，内存存储，无需联网
package noop

import (
	"fmt"
	"sync"

	"github.com/wangbo2295/gadns/core"
	"github.com/wangbo2295/gadns/utils"
)

// Config Noop Provider 配置，与腾讯云 Config 字段一致
type Config struct {
	SecretID  string
	SecretKey string
	Region    string
	Domain    string
}

// Provider Noop CNAMEProvider 实现，行为与腾讯云一致，数据存内存
type Provider struct {
	domain  string
	records map[string][]string // fullDomain → IPs
	mu      sync.RWMutex
}

// NewProvider 创建 Noop Provider
func NewProvider(cfg *Config) *Provider {
	return &Provider{
		domain:  cfg.Domain,
		records: make(map[string][]string),
	}
}

// Create 创建记录
func (p *Provider) Create(fullDomain string, ips []string) (*core.Record, error) {
	for _, ip := range ips {
		if err := utils.ValidateIP(ip); err != nil {
			return nil, fmt.Errorf("invalid IP %q: %w", ip, err)
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.records[fullDomain]; exists {
		return nil, fmt.Errorf("record already exists: %s", fullDomain)
	}

	cp := make([]string, len(ips))
	copy(cp, ips)
	p.records[fullDomain] = cp

	return &core.Record{
		Name:  fullDomain,
		CNAME: utils.GenerateCNAME(fullDomain, p.domain),
		IPs:   ips,
	}, nil
}

// Update 更新记录
func (p *Provider) Update(fullDomain string, ips []string) (*core.Record, error) {
	for _, ip := range ips {
		if err := utils.ValidateIP(ip); err != nil {
			return nil, fmt.Errorf("invalid IP %q: %w", ip, err)
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.records[fullDomain]; !exists {
		return nil, fmt.Errorf("record not found: %s", fullDomain)
	}

	cp := make([]string, len(ips))
	copy(cp, ips)
	p.records[fullDomain] = cp

	return &core.Record{
		Name:  fullDomain,
		CNAME: utils.GenerateCNAME(fullDomain, p.domain),
		IPs:   ips,
	}, nil
}

// Get 获取记录
func (p *Provider) Get(fullDomain string) (*core.Record, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	ips, ok := p.records[fullDomain]
	if !ok {
		return nil, fmt.Errorf("record not found: %s", fullDomain)
	}

	cp := make([]string, len(ips))
	copy(cp, ips)

	return &core.Record{
		Name:  fullDomain,
		CNAME: utils.GenerateCNAME(fullDomain, p.domain),
		IPs:   cp,
	}, nil
}

// List 列出所有记录
func (p *Provider) List() ([]*core.Record, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]*core.Record, 0, len(p.records))
	for name, ips := range p.records {
		cp := make([]string, len(ips))
		copy(cp, ips)
		result = append(result, &core.Record{
			Name:  name,
			CNAME: utils.GenerateCNAME(name, p.domain),
			IPs:   cp,
		})
	}
	return result, nil
}

// Delete 删除记录
func (p *Provider) Delete(fullDomain string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.records[fullDomain]; !ok {
		return fmt.Errorf("record not found: %s", fullDomain)
	}

	delete(p.records, fullDomain)
	return nil
}
