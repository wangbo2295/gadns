// provider/tencent/client_test.go
package tencent

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(&Config{
		SecretID:  "test_id",
		SecretKey: "test_key",
		Region:    "ap-guangzhou",
		Domain:    "example.com",
	})

	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	if client == nil {
		t.Error("NewClient() returned nil")
	}

	if client.domain != "example.com" {
		t.Errorf("domain = %v, want example.com", client.domain)
	}
}

func TestFormatRecordLineRemark(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		carrier  string
		expected string
	}{
		{
			name:     "region only",
			region:   "北京",
			carrier:  "",
			expected: "region:北京",
		},
		{
			name:     "carrier only",
			region:   "",
			carrier:  "电信",
			expected: "carrier:电信",
		},
		{
			name:     "both region and carrier",
			region:   "北京",
			carrier:  "电信",
			expected: "region:北京,carrier:电信",
		},
		{
			name:     "empty",
			region:   "",
			carrier:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatRecordLineRemark(tt.region, tt.carrier)
			if result != tt.expected {
				t.Errorf("FormatRecordLineRemark() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseRecordLineFromRemark(t *testing.T) {
	tests := []struct {
		name          string
		remark        string
		expectedRegion string
		expectedCarrier string
	}{
		{
			name:          "region only",
			remark:        "region:北京",
			expectedRegion: "北京",
			expectedCarrier: "",
		},
		{
			name:          "carrier only",
			remark:        "carrier:电信",
			expectedRegion: "",
			expectedCarrier: "电信",
		},
		{
			name:          "both region and carrier",
			remark:        "region:北京,carrier:电信",
			expectedRegion: "北京",
			expectedCarrier: "电信",
		},
		{
			name:          "with extra spaces",
			remark:        "region:上海, carrier:联通",
			expectedRegion: "上海",
			expectedCarrier: "联通",
		},
		{
			name:          "empty",
			remark:        "",
			expectedRegion: "",
			expectedCarrier: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			region, carrier := ParseRecordLineFromRemark(tt.remark)
			if region != tt.expectedRegion {
				t.Errorf("ParseRecordLineFromRemark() region = %v, want %v", region, tt.expectedRegion)
			}
			if carrier != tt.expectedCarrier {
				t.Errorf("ParseRecordLineFromRemark() carrier = %v, want %v", carrier, tt.expectedCarrier)
			}
		})
	}
}

func TestRecordLineConstants(t *testing.T) {
	tests := []struct {
		name   string
		value  string
	}{
		{"Default", RecordLineDefault},
		{"Telecom", RecordLineTelecom},
		{"Unicom", RecordLineUnicom},
		{"Mobile", RecordLineMobile},
		{"Oversea", RecordLineOversea},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("RecordLine constant %s is empty", tt.name)
			}
		})
	}
}
