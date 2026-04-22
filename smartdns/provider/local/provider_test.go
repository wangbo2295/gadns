// provider/local/provider_test.go
package local_test

import (
	"path/filepath"
	"testing"

	"github.com/yourusername/smartdns/smartdns"
	"github.com/yourusername/smartdns/smartdns/provider/local"
)

func TestLocalProviderCreate(t *testing.T) {
	tmpDir := t.TempDir()
	hostsPath := filepath.Join(tmpDir, "hosts")
	storagePath := filepath.Join(tmpDir, "records.json")

	provider := local.NewProvider(&local.Config{
		HostsPath:   hostsPath,
		StoragePath: storagePath,
		Domain:      "smartdns.local",
	})

	record, err := provider.Create("app", []string{"1.1.1.1", "1.1.1.2"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if record.Name != "app" {
		t.Errorf("Name = %v, want app", record.Name)
	}
	if record.CNAME != "app.smartdns.local" {
		t.Errorf("CNAME = %v, want app.smartdns.local", record.CNAME)
	}
}

func TestLocalProviderGet(t *testing.T) {
	tmpDir := t.TempDir()
	hostsPath := filepath.Join(tmpDir, "hosts")
	storagePath := filepath.Join(tmpDir, "records.json")

	provider := local.NewProvider(&local.Config{
		HostsPath:   hostsPath,
		StoragePath: storagePath,
		Domain:      "smartdns.local",
	})

	// 创建记录
	provider.Create("app", []string{"1.1.1.1"})

	// 获取记录
	record, err := provider.Get("app")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if record.Name != "app" {
		t.Errorf("Name = %v, want app", record.Name)
	}
}

func TestLocalProviderList(t *testing.T) {
	tmpDir := t.TempDir()
	hostsPath := filepath.Join(tmpDir, "hosts")
	storagePath := filepath.Join(tmpDir, "records.json")

	provider := local.NewProvider(&local.Config{
		HostsPath:   hostsPath,
		StoragePath: storagePath,
		Domain:      "smartdns.local",
	})

	// 创建多个记录
	provider.Create("app1", []string{"1.1.1.1"})
	provider.Create("app2", []string{"2.2.2.2"})

	// 列出记录
	records, err := provider.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(records) != 2 {
		t.Errorf("List() returned %d records, want 2", len(records))
	}
}

func TestLocalProviderUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	hostsPath := filepath.Join(tmpDir, "hosts")
	storagePath := filepath.Join(tmpDir, "records.json")

	provider := local.NewProvider(&local.Config{
		HostsPath:   hostsPath,
		StoragePath: storagePath,
		Domain:      "smartdns.local",
	})

	// 创建记录
	provider.Create("app", []string{"1.1.1.1"})

	// 更新记录
	record, err := provider.Update("app", []string{"2.2.2.2", "2.2.2.3"})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if len(record.IPs) != 2 {
		t.Errorf("IPs length = %d, want 2", len(record.IPs))
	}
}

func TestLocalProviderDelete(t *testing.T) {
	tmpDir := t.TempDir()
	hostsPath := filepath.Join(tmpDir, "hosts")
	storagePath := filepath.Join(tmpDir, "records.json")

	provider := local.NewProvider(&local.Config{
		HostsPath:   hostsPath,
		StoragePath: storagePath,
		Domain:      "smartdns.local",
	})

	// 创建记录
	provider.Create("app", []string{"1.1.1.1"})

	// 删除记录
	if err := provider.Delete("app"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// 验证已删除
	_, err := provider.Get("app")
	if err == nil {
		t.Error("Expected error after delete")
	}
}

func TestLocalProviderCreateDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	hostsPath := filepath.Join(tmpDir, "hosts")
	storagePath := filepath.Join(tmpDir, "records.json")

	provider := local.NewProvider(&local.Config{
		HostsPath:   hostsPath,
		StoragePath: storagePath,
		Domain:      "smartdns.local",
	})

	// 创建记录
	provider.Create("app", []string{"1.1.1.1"})

	// 尝试创建重复记录
	_, err := provider.Create("app", []string{"2.2.2.2"})
	if err == nil {
		t.Error("Expected error when creating duplicate record")
	}
}

// TestLocalProviderImplementsInterface 验证 Provider 实现了 SmartDNS 接口
func TestLocalProviderImplementsInterface(t *testing.T) {
	var _ smartdns.SmartDNS = (*local.Provider)(nil)
}
