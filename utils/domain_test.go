package utils

import (
	"strings"
	"testing"
)

func TestSubDomain(t *testing.T) {
	tests := []struct{ full, zone, want string }{
		{"app.example.com", "example.com", "app"},
		{"app.doerhh.cn", "doerhh.cn", "app"},
		{"example.com", "example.com", "@"},
	}
	for _, tt := range tests {
		got := SubDomain(tt.full, tt.zone)
		if got != tt.want {
			t.Errorf("SubDomain(%q, %q) = %q, want %q", tt.full, tt.zone, got, tt.want)
		}
	}
}

func TestFullDomain(t *testing.T) {
	tests := []struct{ sub, zone, want string }{
		{"app", "example.com", "app.example.com"},
		{"@", "example.com", "example.com"},
	}
	for _, tt := range tests {
		got := FullDomain(tt.sub, tt.zone)
		if got != tt.want {
			t.Errorf("FullDomain(%q, %q) = %q, want %q", tt.sub, tt.zone, got, tt.want)
		}
	}
}

func TestGenerateCNAME(t *testing.T) {
	tests := []struct{ full, zone string }{
		{"app.example.com", "example.com"},
		{"ga.doerhh.cn", "doerhh.cn"},
	}
	for _, tt := range tests {
		cname := GenerateCNAME(tt.full, tt.zone)
		sub := SubDomain(tt.full, tt.zone)
		if !strings.HasPrefix(cname, sub+"-") {
			t.Errorf("GenerateCNAME(%q, %q) = %q, want prefix %q-", tt.full, tt.zone, cname, sub)
		}
		if !strings.HasSuffix(cname, "."+tt.zone) {
			t.Errorf("GenerateCNAME(%q, %q) = %q, want suffix .%q", tt.full, tt.zone, cname, tt.zone)
		}
		// 确定性
		if GenerateCNAME(tt.full, tt.zone) != cname {
			t.Errorf("GenerateCNAME not deterministic")
		}
	}
}
