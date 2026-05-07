// provider/tencent/provider.go
package tencent

import (
	"context"
	"fmt"

	"github.com/wangbo2295/gadns/core"
	"github.com/wangbo2295/gadns/utils"
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

// GenerateCNAME 生成 CNAME（委托给 utils）
func (p *Provider) GenerateCNAME(fullDomain string) string {
	return utils.GenerateCNAME(fullDomain, p.config.Domain)
}

// Create 创建DNS记录（负载均衡模式：每个IP一条A记录）
func (p *Provider) Create(fullDomain string, ips []string) (*core.Record, error) {
	ctx := context.Background()

	for _, ip := range ips {
		if err := utils.ValidateIP(ip); err != nil {
			return nil, fmt.Errorf("invalid IP %q: %w", ip, err)
		}
	}

	cname := p.GenerateCNAME(fullDomain)
	sub := utils.SubDomain(cname, p.config.Domain) // A 记录的子域名与 CNAME 一致

	if err := p.createRecords(ctx, sub, ips); err != nil {
		return nil, err
	}

	return &core.Record{
		Name:  fullDomain,
		CNAME: cname,
		IPs:   ips,
	}, nil
}

// createRecords 创建负载均衡记录，失败时回滚已创建的记录
func (p *Provider) createRecords(ctx context.Context, sub string, ips []string) error {
	if len(ips) == 1 {
		_, err := p.client.CreateRecord(ctx, &CreateRecordRequest{
			Domain:     p.config.Domain,
			SubDomain:  sub,
			RecordType: "A",
			RecordLine: "默认",
			Value:      ips[0],
			TTL:        600,
		})
		return err
	}

	weight := uint64(100 / len(ips))
	if weight < 1 {
		weight = 1
	}

	var created []string
	for _, ip := range ips {
		resp, err := p.client.CreateRecord(ctx, &CreateRecordRequest{
			Domain:     p.config.Domain,
			SubDomain:  sub,
			RecordType: "A",
			RecordLine: "默认",
			Value:      ip,
			TTL:        600,
			Weight:     weight,
		})
		if err != nil {
			for _, id := range created {
				p.client.DeleteRecord(ctx, &DeleteRecordRequest{
					Domain:   p.config.Domain,
					RecordID: id,
				})
			}
			return fmt.Errorf("failed to create record for %s: %w", ip, err)
		}
		created = append(created, resp.RecordID)
	}

	return nil
}

// Update 更新DNS记录
func (p *Provider) Update(fullDomain string, ips []string) (*core.Record, error) {
	ctx := context.Background()
	cname := p.GenerateCNAME(fullDomain)
	sub := utils.SubDomain(cname, p.config.Domain)

	recordList, err := p.client.DescribeRecordList(ctx, &DescribeRecordListRequest{
		Subdomain: sub,
		Limit:     100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe records: %w", err)
	}

	for _, record := range recordList.Records {
		if err := p.client.DeleteRecord(ctx, &DeleteRecordRequest{
			Domain:   p.config.Domain,
			RecordID: record.ID,
		}); err != nil {
			return nil, fmt.Errorf("failed to delete old record %s: %w", record.ID, err)
		}
	}

	return p.Create(fullDomain, ips)
}

// Get 获取DNS记录
func (p *Provider) Get(fullDomain string) (*core.Record, error) {
	ctx := context.Background()
	cname := p.GenerateCNAME(fullDomain)
	sub := utils.SubDomain(cname, p.config.Domain)

	recordList, err := p.client.DescribeRecordList(ctx, &DescribeRecordListRequest{
		Subdomain: sub,
		Limit:     100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get record: %w", err)
	}

	if len(recordList.Records) == 0 {
		return nil, fmt.Errorf("record not found: %s", fullDomain)
	}

	ipMap := make(map[string]bool)
	for _, r := range recordList.Records {
		if r.Type == "A" {
			ipMap[r.Value] = true
		}
	}

	ips := make([]string, 0, len(ipMap))
	for ip := range ipMap {
		ips = append(ips, ip)
	}

	return &core.Record{
		Name:  fullDomain,
		CNAME: cname,
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

	recordMap := make(map[string]*core.Record)
	for _, r := range recordList.Records {
		if r.Type != "A" {
			continue
		}

		fd := utils.FullDomain(r.SubDomain, p.config.Domain)
		if existing, ok := recordMap[fd]; ok {
			if !contains(existing.IPs, r.Value) {
				existing.IPs = append(existing.IPs, r.Value)
			}
		} else {
			recordMap[fd] = &core.Record{
				Name:  fd,
				CNAME: fd, // CNAME 即为 A 记录自身的域名
				IPs:   []string{r.Value},
			}
		}
	}

	result := make([]*core.Record, 0, len(recordMap))
	for _, r := range recordMap {
		result = append(result, r)
	}

	return result, nil
}

// Delete 删除DNS记录
func (p *Provider) Delete(fullDomain string) error {
	ctx := context.Background()
	cname := p.GenerateCNAME(fullDomain)
	sub := utils.SubDomain(cname, p.config.Domain)

	recordList, err := p.client.DescribeRecordList(ctx, &DescribeRecordListRequest{
		Subdomain: sub,
		Limit:     100,
	})
	if err != nil {
		return fmt.Errorf("failed to get records: %w", err)
	}

	if len(recordList.Records) == 0 {
		return fmt.Errorf("record not found: %s", fullDomain)
	}

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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
