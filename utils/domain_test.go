package utils

import (
	"strings"
	"testing"
)

func TestGenerateCNAME(t *testing.T) {
	tests := []struct{ full, zone string }{
		{"app.example.com", "example.com"},
		{"app.doerhh.cn", "doerhh.cn"},
		{"web.doerhh.cn", "doerhh.cn"},
	}

	for _, tt := range tests {
		cname := GenerateCNAME(tt.full, tt.zone)
		// CNAME 应以 zoneDomain 结尾
		if !strings.HasSuffix(cname, "."+tt.zone) {
			t.Errorf("GenerateCNAME(%q, %q) = %q, want suffix .%q", tt.full, tt.zone, cname, tt.zone)
		}
		// 确定性
		if GenerateCNAME(tt.full, tt.zone) != cname {
			t.Errorf("GenerateCNAME not deterministic")
		}
	}

	// 不同 fullDomain 应产生不同 CNAME
	a := GenerateCNAME("app.example.com", "example.com")
	b := GenerateCNAME("web.example.com", "example.com")
	if a == b {
		t.Errorf("different domains should produce different CNAMEs")
	}

	// 多级子域名也能正确生成
	cname := GenerateCNAME("a.b.example.com", "example.com")
	if !strings.HasPrefix(cname, "a-b-example-com-") {
		t.Errorf("multi-level domain CNAME = %q", cname)
	}
}
