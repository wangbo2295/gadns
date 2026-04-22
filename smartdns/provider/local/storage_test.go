// provider/local/storage_test.go
package local

import (
	"path/filepath"
	"testing"
)

func TestStorageSaveAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	storagePath := filepath.Join(tmpDir, "records.json")

	storage := NewStorage(storagePath)

	record := &StoredRecord{
		Name:      "app",
		CNAME:     "app.smartdns.local",
		IPs:       []string{"1.1.1.1", "1.1.1.2"},
		CurrentIP: "1.1.1.1",
	}

	// 保存记录
	if err := storage.Save("app", record); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// 获取记录
	got, err := storage.Get("app")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.Name != record.Name {
		t.Errorf("Name = %v, want %v", got.Name, record.Name)
	}
	if got.CNAME != record.CNAME {
		t.Errorf("CNAME = %v, want %v", got.CNAME, record.CNAME)
	}
}

func TestStorageList(t *testing.T) {
	tmpDir := t.TempDir()
	storagePath := filepath.Join(tmpDir, "records.json")

	storage := NewStorage(storagePath)

	// 保存多个记录
	storage.Save("app1", &StoredRecord{Name: "app1", CNAME: "app1.smartdns.local"})
	storage.Save("app2", &StoredRecord{Name: "app2", CNAME: "app2.smartdns.local"})

	// 列出所有记录
	records, err := storage.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(records) != 2 {
		t.Errorf("List() returned %d records, want 2", len(records))
	}
}

func TestStorageDelete(t *testing.T) {
	tmpDir := t.TempDir()
	storagePath := filepath.Join(tmpDir, "records.json")

	storage := NewStorage(storagePath)

	// 保存记录
	storage.Save("app", &StoredRecord{Name: "app", CNAME: "app.smartdns.local"})

	// 删除记录
	if err := storage.Delete("app"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// 验证已删除
	_, err := storage.Get("app")
	if err == nil {
		t.Error("Expected error after delete")
	}
}
