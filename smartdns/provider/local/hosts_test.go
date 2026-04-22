// provider/local/hosts_test.go
package local

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHostsManagerAddEntry(t *testing.T) {
	tmpDir := t.TempDir()
	hostsPath := filepath.Join(tmpDir, "hosts")

	hm := NewHostsManager(hostsPath)

	// 添加记录
	if err := hm.AddEntry("127.0.0.1", "app.smartdns.local"); err != nil {
		t.Fatalf("AddEntry() error = %v", err)
	}

	// 读取并验证
	content, err := os.ReadFile(hostsPath)
	if err != nil {
		t.Fatalf("Failed to read hosts file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "127.0.0.1") || !strings.Contains(contentStr, "app.smartdns.local") {
		t.Errorf("Hosts file does not contain expected entry")
	}
}

func TestHostsManagerUpdateEntry(t *testing.T) {
	tmpDir := t.TempDir()
	hostsPath := filepath.Join(tmpDir, "hosts")

	hm := NewHostsManager(hostsPath)

	// 添加初始记录
	hm.AddEntry("1.1.1.1", "app.smartdns.local")

	// 更新记录
	if err := hm.UpdateEntry("2.2.2.2", "app.smartdns.local"); err != nil {
		t.Fatalf("UpdateEntry() error = %v", err)
	}

	// 读取并验证
	content, _ := os.ReadFile(hostsPath)
	contentStr := string(content)

	if strings.Contains(contentStr, "1.1.1.1") {
		t.Error("Old IP still present in hosts file")
	}
	if !strings.Contains(contentStr, "2.2.2.2") {
		t.Error("New IP not found in hosts file")
	}
}

func TestHostsManagerRemoveEntry(t *testing.T) {
	tmpDir := t.TempDir()
	hostsPath := filepath.Join(tmpDir, "hosts")

	hm := NewHostsManager(hostsPath)

	// 添加记录
	hm.AddEntry("1.1.1.1", "app.smartdns.local")

	// 删除记录
	if err := hm.RemoveEntry("app.smartdns.local"); err != nil {
		t.Fatalf("RemoveEntry() error = %v", err)
	}

	// 读取并验证
	content, _ := os.ReadFile(hostsPath)
	contentStr := string(content)

	if strings.Contains(contentStr, "app.smartdns.local") {
		t.Error("Entry still present in hosts file")
	}
}
