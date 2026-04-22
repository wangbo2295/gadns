// provider/tencent/client.go
package tencent

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	credential := tencentcommon.NewCredential(
		cfg.SecretID,
		cfg.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"

	client, err := dnspod.NewClient(credential, cfg.Region, cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to create dnspod client: %w", err)
	}

	return &Client{
		domain: cfg.Domain,
		client: client,
	}, nil
}

// CreateRecord 创建DNS记录
func (c *Client) CreateRecord(ctx context.Context, req *CreateRecordRequest) (*CreateRecordResponse, error) {
	dnspodReq := dnspod.NewCreateRecordRequest()

	// 设置域名
	domain := c.domain
	if req.Domain != "" {
		domain = req.Domain
	}
	dnspodReq.Domain = &domain

	// 设置子域名
	dnspodReq.SubDomain = &req.SubDomain

	// 设置记录类型（A记录）
	dnspodReq.RecordType = &req.RecordType

	// 设置记录线路（智能调度）
	if req.RecordLine != "" {
		dnspodReq.RecordLine = &req.RecordLine
	}
	if req.RecordLineID != "" {
		dnspodReq.RecordLineId = &req.RecordLineID
	}

	// 设置记录值（IP）
	dnspodReq.Value = &req.Value

	// 设置TTL
	dnspodReq.TTL = &req.TTL

	// 设置权重（用于负载均衡）
	if req.Weight > 0 {
		dnspodReq.Weight = &req.Weight
	}

	// 设置备注
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

// ModifyRecord 修改DNS记录
func (c *Client) ModifyRecord(ctx context.Context, req *ModifyRecordRequest) error {
	dnspodReq := dnspod.NewModifyRecordRequest()

	// 设置域名
	domain := c.domain
	if req.Domain != "" {
		domain = req.Domain
	}
	dnspodReq.Domain = &domain

	// 设置记录ID
	recordID, err := strconv.ParseUint(req.RecordID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid RecordID: %w", err)
	}
	dnspodReq.RecordId = &recordID

	// 设置子域名
	dnspodReq.SubDomain = &req.SubDomain

	// 设置记录类型
	dnspodReq.RecordType = &req.RecordType

	// 设置记录线路
	if req.RecordLine != "" {
		dnspodReq.RecordLine = &req.RecordLine
	}
	if req.RecordLineID != "" {
		dnspodReq.RecordLineId = &req.RecordLineID
	}

	// 设置记录值
	dnspodReq.Value = &req.Value

	// 设置TTL
	if req.TTL > 0 {
		dnspodReq.TTL = &req.TTL
	}

	// 设置权重
	if req.Weight > 0 {
		dnspodReq.Weight = &req.Weight
	}

	// 设置备注
	if req.Remark != "" {
		dnspodReq.Remark = &req.Remark
	}

	_, err = c.client.ModifyRecord(dnspodReq)
	if err != nil {
		return fmt.Errorf("ModifyRecord API failed: %w", err)
	}

	return nil
}

// DeleteRecord 删除DNS记录
func (c *Client) DeleteRecord(ctx context.Context, req *DeleteRecordRequest) error {
	dnspodReq := dnspod.NewDeleteRecordRequest()

	// 设置域名
	domain := c.domain
	if req.Domain != "" {
		domain = req.Domain
	}
	dnspodReq.Domain = &domain

	// 设置记录ID
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

	// 设置域名
	domain := c.domain
	if req.Domain != "" {
		domain = req.Domain
	}
	dnspodReq.Domain = &domain

	// 设置子域名过滤
	if req.Subdomain != "" {
		dnspodReq.Subdomain = &req.Subdomain
	}

	// 设置记录类型过滤
	if req.RecordType != "" {
		dnspodReq.RecordType = &req.RecordType
	}

	// 设置记录线路过滤
	if req.RecordLine != "" {
		dnspodReq.RecordLine = &req.RecordLine
	}

	// 设置分页
	limit := uint64(req.Limit)
	if limit == 0 {
		limit = 100 // 默认100条
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
		record := &DNSRecord{
			SubDomain: *r.Name,
		}

		// 记录ID
		if r.RecordId != nil {
			record.ID = strconv.FormatUint(*r.RecordId, 10)
		}

		// 基本字段
		if r.Name != nil {
			record.Name = *r.Name
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

		// 记录线路信息 (使用 Line 和 LineId)
		if r.Line != nil {
			record.RecordLine = *r.Line
		}
		if r.LineId != nil {
			record.RecordLineID = *r.LineId
		}

		// 备注信息
		if r.Remark != nil {
			record.Remark = *r.Remark
		}

		// 更新时间
		if r.UpdatedOn != nil {
			record.UpdatedOn = *r.UpdatedOn
		}

		records = append(records, record)
	}

	totalCount := 0
	if response.Response.RecordCountInfo != nil && response.Response.RecordCountInfo.TotalCount != nil {
		totalCount = int(*response.Response.RecordCountInfo.TotalCount)
	} else {
		totalCount = len(records)
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

// DNSRecord DNS记录
type DNSRecord struct {
	ID           string `json:"id"`            // 记录ID
	Name         string `json:"name"`          // 记录名称
	Type         string `json:"type"`          // 记录类型
	Value        string `json:"value"`         // 记录值
	TTL          int    `json:"ttl"`           // TTL
	Weight       int    `json:"weight"`        // 权重
	SubDomain    string `json:"sub_domain"`    // 子域名
	RecordLine   string `json:"record_line"`   // 记录线路
	RecordLineID string `json:"record_line_id"`// 记录线路ID
	Remark       string `json:"remark"`        // 备注
	UpdatedOn    string `json:"updated_on"`    // 更新时间
}

// CreateRecordRequest 创建记录请求
type CreateRecordRequest struct {
	Domain       string // 域名（默认使用配置的域名）
	SubDomain    string // 子域名（必填）
	RecordType   string // 记录类型：A、CNAME、AAAA等
	RecordLine   string // 记录线路：默认、电信、联通等
	RecordLineID string // 记录线路ID
	Value        string // 记录值（IP地址或域名）
	TTL          uint64 // TTL（秒）
	Weight       uint64 // 权重（1-100）
	Remark       string // 备注
}

// CreateRecordResponse 创建记录响应
type CreateRecordResponse struct {
	RecordID string // 记录ID
	Name     string // 记录名称
	Type     string // 记录类型
	Value    string // 记录值
}

// ModifyRecordRequest 修改记录请求
type ModifyRecordRequest struct {
	Domain       string // 域名
	RecordID     string // 记录ID（必填）
	SubDomain    string // 子域名
	RecordType   string // 记录类型
	RecordLine   string // 记录线路
	RecordLineID string // 记录线路ID
	Value        string // 记录值
	TTL          uint64 // TTL
	Weight       uint64 // 权重
	Remark       string // 备注
}

// DeleteRecordRequest 删除记录请求
type DeleteRecordRequest struct {
	Domain   string // 域名
	RecordID string // 记录ID（必填）
}

// DescribeRecordListRequest 查询记录列表请求
type DescribeRecordListRequest struct {
	Domain      string // 域名
	Subdomain   string // 子域名（过滤条件）
	RecordType  string // 记录类型（过滤条件）
	RecordLine  string // 记录线路（过滤条件）
	Limit       int    // 分页限制
	Offset      int    // 分页偏移
}

// DescribeRecordListResponse 查询记录列表响应
type DescribeRecordListResponse struct {
	Records    []*DNSRecord `json:"records"`     // 记录列表
	TotalCount int          `json:"total_count"` // 总数量
}

// RecordLine 记录线路枚举
const (
	RecordLineDefault = "默认"      // 默认线路
	RecordLineTelecom = "电信"      // 电信
	RecordLineUnicom  = "联通"      // 联通
	RecordLineMobile  = "移动"      // 移动
	RecordLineOversea = "境外"      // 境外
	RecordLineBaidu   = "百度"      // 百度
	RecordLineGoogle  = "谷歌"      // 谷歌
	RecordLineKnown   = "知道"      // 知道
	RecordLineBing    = "必应"      // 必应
	RecordLineSougo   = "搜狗"      // 搜狗
	RecordLineOther   = "其他"      // 其他
)

// ProvinceRecordLines 省份线路枚举
var ProvinceRecordLines = []string{
	"北京", "上海", "天津", "重庆", "广东", "江苏", "浙江", "四川", "湖北", "湖南",
	"福建", "山东", "河南", "河北", "山西", "陕西", "辽宁", "吉林", "黑龙江", "安徽",
	"江西", "广西", "海南", "云南", "贵州", "甘肃", "青海", "宁夏", "新疆", "内蒙古",
}

// ParseRecordLineFromRemark 从备注中解析线路信息
// 格式：region:北京 或 carrier:电信 或 region:北京,carrier:电信
func ParseRecordLineFromRemark(remark string) (region, carrier string) {
	parts := strings.Split(remark, ",")
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), ":", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "region":
			region = kv[1]
		case "carrier":
			carrier = kv[1]
		}
	}
	return
}

// FormatRecordLineRemark 格式化线路备注
func FormatRecordLineRemark(region, carrier string) string {
	var parts []string
	if region != "" {
		parts = append(parts, fmt.Sprintf("region:%s", region))
	}
	if carrier != "" {
		parts = append(parts, fmt.Sprintf("carrier:%s", carrier))
	}
	return strings.Join(parts, ",")
}
