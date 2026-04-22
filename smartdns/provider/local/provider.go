// provider/local/provider.go
package local

import (
	"fmt"
	"time"

	"github.com/yourusername/smartdns/smartdns"
	"github.com/yourusername/smartdns/smartdns/iputil"
)

// Provider 本地DNS Provider
type Provider struct {
	config   *Config
	storage  *Storage
	hostsMgr *HostsManager
}

// NewProvider 创建本地Provider实例
func NewProvider(cfg *Config) *Provider {
	return &Provider{
		config:   cfg,
		storage:  NewStorage(cfg.StoragePath),
		hostsMgr: NewHostsManager(cfg.HostsPath),
	}
}

// Create 创建DNS记录
func (p *Provider) Create(name string, ips []string) (*smartdns.Record, error) {
	// 检查记录是否已存在
	_, err := p.storage.Get(name)
	if err == nil {
		return nil, fmt.Errorf("record already exists: %s", name)
	}

	// 解析IP
	parsedIPs, err := iputil.ParseIPs(ips)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IPs: %w", err)
	}

	// 选择当前IP（简单轮询：取第一个）
	currentIP := parsedIPs[0]

	cname := fmt.Sprintf("%s.%s", name, p.config.Domain)

	// 存储记录
	storedRecord := &StoredRecord{
		Name:      name,
		CNAME:     cname,
		IPs:       ips,
		CurrentIP: currentIP,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	if err := p.storage.Save(name, storedRecord); err != nil {
		return nil, err
	}

	// 更新hosts文件
	if err := p.hostsMgr.AddEntry(currentIP, cname); err != nil {
		return nil, fmt.Errorf("failed to update hosts file: %w", err)
	}

	return &smartdns.Record{
		Name:  name,
		CNAME: cname,
		IPs:   ips,
	}, nil
}

// Update 更新DNS记录
func (p *Provider) Update(name string, ips []string) (*smartdns.Record, error) {
	// 检查记录是否存在
	_, err := p.storage.Get(name)
	if err != nil {
		return nil, fmt.Errorf("record not found: %s", name)
	}

	// 解析IP
	parsedIPs, err := iputil.ParseIPs(ips)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IPs: %w", err)
	}

	// 选择当前IP（简单轮询：取第一个）
	currentIP := parsedIPs[0]

	cname := fmt.Sprintf("%s.%s", name, p.config.Domain)

	// 更新存储
	storedRecord := &StoredRecord{
		Name:      name,
		CNAME:     cname,
		IPs:       ips,
		CurrentIP: currentIP,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	if err := p.storage.Save(name, storedRecord); err != nil {
		return nil, err
	}

	// 更新hosts文件
	if err := p.hostsMgr.UpdateEntry(currentIP, cname); err != nil {
		return nil, fmt.Errorf("failed to update hosts file: %w", err)
	}

	return &smartdns.Record{
		Name:  name,
		CNAME: cname,
		IPs:   ips,
	}, nil
}

// Get 获取DNS记录
func (p *Provider) Get(name string) (*smartdns.Record, error) {
	storedRecord, err := p.storage.Get(name)
	if err != nil {
		return nil, err
	}

	return &smartdns.Record{
		Name:  storedRecord.Name,
		CNAME: storedRecord.CNAME,
		IPs:   storedRecord.IPs,
	}, nil
}

// List 列出所有DNS记录
func (p *Provider) List() ([]*smartdns.Record, error) {
	storedRecords, err := p.storage.List()
	if err != nil {
		return nil, err
	}

	var records []*smartdns.Record
	for _, sr := range storedRecords {
		records = append(records, &smartdns.Record{
			Name:  sr.Name,
			CNAME: sr.CNAME,
			IPs:   sr.IPs,
		})
	}

	return records, nil
}

// Delete 删除DNS记录
func (p *Provider) Delete(name string) error {
	// 获取记录以获取CNAME
	storedRecord, err := p.storage.Get(name)
	if err != nil {
		return err
	}

	// 从hosts文件删除
	if err := p.hostsMgr.RemoveEntry(storedRecord.CNAME); err != nil {
		return fmt.Errorf("failed to remove from hosts file: %w", err)
	}

	// 从存储删除
	return p.storage.Delete(name)
}
