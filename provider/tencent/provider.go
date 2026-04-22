// provider/tencent/provider.go
package tencent

import (
	"context"
	"fmt"
	"strings"

	"github.com/yourusername/smartdns/core"
	"github.com/yourusername/smartdns/iputil"
)

// Provider 腾讯云DNS Provider
type Provider struct {
	config *Config
	client *Client
}

// NewProvider 创建腾讯云Provider实例
func NewProvider(cfg *Config) (*Provider, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Provider{
		config: cfg,
		client: client,
	}, nil
}

// GenerateCNAME 生成CNAME
func (p *Provider) GenerateCNAME(name string) string {
	return fmt.Sprintf("%s.%s", name, p.config.Domain)
}

// SmartRoutingConfig 智能调度配置
type SmartRoutingConfig struct {
	Enabled bool     // 是否启用智能调度
	Regions []string // 地域列表
	Carriers []string // 运营商列表
}

// DefaultSmartRoutingConfig 默认智能调度配置
var DefaultSmartRoutingConfig = SmartRoutingConfig{
	Enabled: true,
	Regions: []string{
		"北京", "上海", "广东", "江苏", "浙江", "四川", "湖北", "福建",
	},
	Carriers: []string{
		RecordLineTelecom, RecordLineUnicom, RecordLineMobile,
	},
}

// Create 创建DNS记录（支持智能调度）
func (p *Provider) Create(name string, ips []string) (*core.Record, error) {
	ctx := context.Background()

	// 解析IP列表
	parsedIPs, err := iputil.ParseIPs(ips)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IPs: %w", err)
	}

	cname := p.GenerateCNAME(name)

	// 检查是否启用智能调度
	config := p.getSmartRoutingConfig()

	if config.Enabled {
		// 智能调度模式：为每个IP创建多条记录
		_, err := p.createSmartRoutingRecords(ctx, name, parsedIPs, config)
		if err != nil {
			return nil, err
		}

		return &core.Record{
			Name:  name,
			CNAME: cname,
			IPs:   ips,
		}, nil
	}

	// 普通模式：创建单条记录（多个IP用换行分隔）
	ipList := strings.Join(parsedIPs, "\n")

	_, err = p.createSingleRecord(ctx, name, ipList, RecordLineDefault, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to create record: %w", err)
	}

	return &core.Record{
		Name:  name,
		CNAME: cname,
		IPs:   ips,
	}, nil
}

// createSmartRoutingRecords 创建智能调度记录
func (p *Provider) createSmartRoutingRecords(ctx context.Context, name string, ips []string, config SmartRoutingConfig) ([]string, error) {
	var recordIDs []string

	// 为每个IP创建多条线路记录
	for ipIndex, ip := range ips {
		// 创建默认线路记录（第一条IP）
		if ipIndex == 0 {
			recordID, err := p.createSingleRecord(ctx, name, ip, RecordLineDefault, "", fmt.Sprintf("primary:%d", ipIndex))
			if err != nil {
				return nil, fmt.Errorf("failed to create default record: %w", err)
			}
			recordIDs = append(recordIDs, recordID)
		}

		// 创建地域线路记录
		for regionIndex, region := range config.Regions {
			recordID, err := p.createSingleRecord(
				ctx, name, ip,
				region, "",
				fmt.Sprintf(FormatRecordLineRemark(region, "")+",index:%d", ipIndex*len(config.Regions)+regionIndex),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create region record %s: %w", region, err)
			}
			recordIDs = append(recordIDs, recordID)
		}

		// 创建运营商线路记录
		for _, carrier := range config.Carriers {
			recordID, err := p.createSingleRecord(
				ctx, name, ip,
				carrier, "",
				fmt.Sprintf(FormatRecordLineRemark("", carrier)+",index:%d", ipIndex),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create carrier record %s: %w", carrier, err)
			}
			recordIDs = append(recordIDs, recordID)
		}
	}

	return recordIDs, nil
}

// createSingleRecord 创建单条记录
func (p *Provider) createSingleRecord(ctx context.Context, subdomain, value, recordLine, recordLineID, remark string) (string, error) {
	req := &CreateRecordRequest{
		Domain:       p.config.Domain,
		SubDomain:    subdomain,
		RecordType:   "A",
		RecordLine:   recordLine,
		RecordLineID: recordLineID,
		Value:        value,
		TTL:          600,
		Weight:       1,
		Remark:       remark,
	}

	resp, err := p.client.CreateRecord(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.RecordID, nil
}

// Update 更新DNS记录
func (p *Provider) Update(name string, ips []string) (*core.Record, error) {
	ctx := context.Background()

	// 解析IP列表
	parsedIPs, err := iputil.ParseIPs(ips)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IPs: %w", err)
	}

	cname := p.GenerateCNAME(name)

	// 获取现有记录列表
	recordList, err := p.client.DescribeRecordList(ctx, &DescribeRecordListRequest{
		Subdomain: name,
		Limit:     100, // 获取所有相关记录
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe records: %w", err)
	}

	// 如果没有现有记录，转为创建
	if len(recordList.Records) == 0 {
		return p.Create(name, ips)
	}

	config := p.getSmartRoutingConfig()

	if config.Enabled && len(recordList.Records) > 1 {
		// 智能调度模式：删除所有旧记录，重新创建
		for _, record := range recordList.Records {
			if err := p.client.DeleteRecord(ctx, &DeleteRecordRequest{
				Domain:   p.config.Domain,
				RecordID: record.ID,
			}); err != nil {
				return nil, fmt.Errorf("failed to delete old record %s: %w", record.ID, err)
			}
		}

		// 重新创建记录
		_, err := p.createSmartRoutingRecords(ctx, name, parsedIPs, config)
		if err != nil {
			return nil, err
		}

		return &core.Record{
			Name:  name,
			CNAME: cname,
			IPs:   ips,
		}, nil
	}

	// 普通模式：更新单条记录
	record := recordList.Records[0]
	ipList := strings.Join(parsedIPs, "\n")

	if err := p.client.ModifyRecord(ctx, &ModifyRecordRequest{
		Domain:       p.config.Domain,
		RecordID:     record.ID,
		SubDomain:    name,
		RecordType:   "A",
		RecordLine:   record.RecordLine, // SDK 中使用 Line 字段
		RecordLineID: record.RecordLineID, // SDK 中使用 LineId 字段
		Value:        ipList,
		TTL:          uint64(record.TTL),
		Weight:       uint64(record.Weight),
	}); err != nil {
		return nil, fmt.Errorf("failed to modify record: %w", err)
	}

	return &core.Record{
		Name:  name,
		CNAME: cname,
		IPs:   ips,
	}, nil
}

// Get 获取DNS记录
func (p *Provider) Get(name string) (*core.Record, error) {
	ctx := context.Background()

	recordList, err := p.client.DescribeRecordList(ctx, &DescribeRecordListRequest{
		Subdomain: name,
		Limit:     100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get record: %w", err)
	}

	if len(recordList.Records) == 0 {
		return nil, fmt.Errorf("record not found: %s", name)
	}

	// 收集所有IP值（去重）
	ipMap := make(map[string]bool)
	for _, record := range recordList.Records {
		if record.Type == "A" {
			// 处理多值记录（换行分隔）
			values := strings.Split(record.Value, "\n")
			for _, v := range values {
				v = strings.TrimSpace(v)
				if v != "" {
					ipMap[v] = true
				}
			}
		}
	}

	ips := make([]string, 0, len(ipMap))
	for ip := range ipMap {
		ips = append(ips, ip)
	}

	return &core.Record{
		Name:  name,
		CNAME: p.GenerateCNAME(name),
		IPs:   ips,
	}, nil
}

// List 列出所有DNS记录
func (p *Provider) List() ([]*core.Record, error) {
	ctx := context.Background()

	recordList, err := p.client.DescribeRecordList(ctx, &DescribeRecordListRequest{
		Limit: 1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list records: %w", err)
	}

	// 按子域名分组记录
	recordMap := make(map[string]*core.Record)
	for _, record := range recordList.Records {
		if record.Type != "A" {
			continue
		}

		if existing, ok := recordMap[record.SubDomain]; ok {
			// 添加新的IP值
			values := strings.Split(record.Value, "\n")
			for _, v := range values {
				v = strings.TrimSpace(v)
				if v != "" && !contains(existing.IPs, v) {
					existing.IPs = append(existing.IPs, v)
				}
			}
		} else {
			// 创建新记录
			values := strings.Split(record.Value, "\n")
			ips := make([]string, 0, len(values))
			for _, v := range values {
				v = strings.TrimSpace(v)
				if v != "" {
					ips = append(ips, v)
				}
			}

			recordMap[record.SubDomain] = &core.Record{
				Name:  record.SubDomain,
				CNAME: p.GenerateCNAME(record.SubDomain),
				IPs:   ips,
			}
		}
	}

	// 转换为数组
	result := make([]*core.Record, 0, len(recordMap))
	for _, record := range recordMap {
		result = append(result, record)
	}

	return result, nil
}

// Delete 删除DNS记录
func (p *Provider) Delete(name string) error {
	ctx := context.Background()

	recordList, err := p.client.DescribeRecordList(ctx, &DescribeRecordListRequest{
		Subdomain: name,
		Limit:     100,
	})
	if err != nil {
		return fmt.Errorf("failed to get records: %w", err)
	}

	if len(recordList.Records) == 0 {
		return fmt.Errorf("record not found: %s", name)
	}

	// 删除所有匹配的记录
	for _, record := range recordList.Records {
		if err := p.client.DeleteRecord(ctx, &DeleteRecordRequest{
			Domain:   p.config.Domain,
			RecordID: record.ID,
		}); err != nil {
			return fmt.Errorf("failed to delete record %s: %w", record.ID, err)
		}
	}

	return nil
}

// getSmartRoutingConfig 获取智能调度配置
func (p *Provider) getSmartRoutingConfig() SmartRoutingConfig {
	// TODO: 从配置文件或环境变量读取
	// 目前返回默认配置
	return DefaultSmartRoutingConfig
}

// contains 检查字符串切片是否包含某元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
