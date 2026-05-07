// provider/tencent/client.go
package tencent

import (
	"context"
	"fmt"
	"strconv"

	tencentcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

// Client 腾讯云DNS客户端
type Client struct {
	domain string
	client *dnspod.Client
}

// NewClient 创建腾讯云DNS客户端
func NewClient(cfg *Config) (*Client, error) {
	credential := tencentcommon.NewCredential(cfg.SecretID, cfg.SecretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"

	client, err := dnspod.NewClient(credential, cfg.Region, cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to create dnspod client: %w", err)
	}

	return &Client{domain: cfg.Domain, client: client}, nil
}

// CreateRecord 创建DNS记录
func (c *Client) CreateRecord(ctx context.Context, req *CreateRecordRequest) (*CreateRecordResponse, error) {
	dnspodReq := dnspod.NewCreateRecordRequest()

	domain := c.domain
	if req.Domain != "" {
		domain = req.Domain
	}
	dnspodReq.Domain = &domain
	dnspodReq.SubDomain = &req.SubDomain
	dnspodReq.RecordType = &req.RecordType
	dnspodReq.Value = &req.Value
	dnspodReq.TTL = &req.TTL

	if req.RecordLine != "" {
		dnspodReq.RecordLine = &req.RecordLine
	}
	if req.Weight > 0 {
		dnspodReq.Weight = &req.Weight
	}
	if req.Remark != "" {
		dnspodReq.Remark = &req.Remark
	}

	response, err := c.client.CreateRecord(dnspodReq)
	if err != nil {
		return nil, fmt.Errorf("CreateRecord API failed: %w", err)
	}

	if response.Response.RecordId == nil {
		return nil, fmt.Errorf("CreateRecord response missing RecordId")
	}

	return &CreateRecordResponse{
		RecordID: strconv.FormatUint(*response.Response.RecordId, 10),
		Name:     req.SubDomain,
		Type:     req.RecordType,
		Value:    req.Value,
	}, nil
}

// DeleteRecord 删除DNS记录
func (c *Client) DeleteRecord(ctx context.Context, req *DeleteRecordRequest) error {
	dnspodReq := dnspod.NewDeleteRecordRequest()

	domain := c.domain
	if req.Domain != "" {
		domain = req.Domain
	}
	dnspodReq.Domain = &domain

	recordID, err := strconv.ParseUint(req.RecordID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid RecordID: %w", err)
	}
	dnspodReq.RecordId = &recordID

	_, err = c.client.DeleteRecord(dnspodReq)
	if err != nil {
		return fmt.Errorf("DeleteRecord API failed: %w", err)
	}

	return nil
}

// DescribeRecordList 查询DNS记录列表
func (c *Client) DescribeRecordList(ctx context.Context, req *DescribeRecordListRequest) (*DescribeRecordListResponse, error) {
	dnspodReq := dnspod.NewDescribeRecordListRequest()

	domain := c.domain
	if req.Domain != "" {
		domain = req.Domain
	}
	dnspodReq.Domain = &domain

	if req.Subdomain != "" {
		dnspodReq.Subdomain = &req.Subdomain
	}
	if req.RecordType != "" {
		dnspodReq.RecordType = &req.RecordType
	}

	limit := uint64(req.Limit)
	if limit == 0 {
		limit = 100
	}
	dnspodReq.Limit = &limit

	if req.Offset > 0 {
		offset := uint64(req.Offset)
		dnspodReq.Offset = &offset
	}

	response, err := c.client.DescribeRecordList(dnspodReq)
	if err != nil {
		return nil, fmt.Errorf("DescribeRecordList API failed: %w", err)
	}

	records := make([]*DNSRecord, 0, len(response.Response.RecordList))
	for _, r := range response.Response.RecordList {
		record := &DNSRecord{}

		if r.RecordId != nil {
			record.ID = strconv.FormatUint(*r.RecordId, 10)
		}
		if r.Name != nil {
			record.Name = *r.Name
			record.SubDomain = *r.Name
		}
		if r.Type != nil {
			record.Type = *r.Type
		}
		if r.Value != nil {
			record.Value = *r.Value
		}
		if r.TTL != nil {
			record.TTL = int(*r.TTL)
		}
		if r.Weight != nil {
			record.Weight = int(*r.Weight)
		}
		if r.UpdatedOn != nil {
			record.UpdatedOn = *r.UpdatedOn
		}

		records = append(records, record)
	}

	totalCount := len(records)
	if response.Response.RecordCountInfo != nil && response.Response.RecordCountInfo.TotalCount != nil {
		totalCount = int(*response.Response.RecordCountInfo.TotalCount)
	}

	return &DescribeRecordListResponse{
		Records:    records,
		TotalCount: totalCount,
	}, nil
}

// GetRecordIDBySubDomain 根据子域名获取记录ID
func (c *Client) GetRecordIDBySubDomain(ctx context.Context, subdomain string) (string, error) {
	resp, err := c.DescribeRecordList(ctx, &DescribeRecordListRequest{
		Subdomain: subdomain,
		Limit:     1,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Records) == 0 {
		return "", fmt.Errorf("record not found: %s", subdomain)
	}

	return resp.Records[0].ID, nil
}

// --- Types ---

// DNSRecord DNS记录
type DNSRecord struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	TTL       int    `json:"ttl"`
	Weight    int    `json:"weight"`
	SubDomain string `json:"sub_domain"`
	UpdatedOn string `json:"updated_on"`
}

// CreateRecordRequest 创建记录请求
type CreateRecordRequest struct {
	Domain     string
	SubDomain  string
	RecordType string
	RecordLine string
	Value      string
	TTL        uint64
	Weight     uint64
	Remark     string
}

// CreateRecordResponse 创建记录响应
type CreateRecordResponse struct {
	RecordID string
	Name     string
	Type     string
	Value    string
}

// DeleteRecordRequest 删除记录请求
type DeleteRecordRequest struct {
	Domain   string
	RecordID string
}

// DescribeRecordListRequest 查询记录列表请求
type DescribeRecordListRequest struct {
	Domain     string
	Subdomain  string
	RecordType string
	Limit      int
	Offset     int
}

// DescribeRecordListResponse 查询记录列表响应
type DescribeRecordListResponse struct {
	Records    []*DNSRecord `json:"records"`
	TotalCount int          `json:"total_count"`
}
