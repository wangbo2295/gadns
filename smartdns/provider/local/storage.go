// provider/local/storage.go
package local

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// StoredRecord 存储的记录
type StoredRecord struct {
	Name      string   `json:"name"`
	CNAME     string   `json:"cname"`
	IPs       []string `json:"ips"`
	CurrentIP string   `json:"current_ip"`
	UpdatedAt string   `json:"updated_at"`
}

// Storage 记录存储
type Storage struct {
	path   string
	mu     sync.RWMutex
	records map[string]*StoredRecord
}

// NewStorage 创建存储实例
func NewStorage(path string) *Storage {
	return &Storage{
		path:    path,
		records: make(map[string]*StoredRecord),
	}
}

// load 从文件加载记录
func (s *Storage) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，初始化空记录
			s.records = make(map[string]*StoredRecord)
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, &s.records); err != nil {
		return err
	}

	return nil
}

// save 保存记录到文件
func (s *Storage) save() error {
	data, err := json.MarshalIndent(s.records, "", "  ")
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}

// Get 获取记录
func (s *Storage) Get(name string) (*StoredRecord, error) {
	if err := s.load(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	record, ok := s.records[name]
	if !ok {
		return nil, fmt.Errorf("record not found: %s", name)
	}

	return record, nil
}

// Save 保存记录
func (s *Storage) Save(name string, record *StoredRecord) error {
	if err := s.load(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.records[name] = record
	return s.save()
}

// List 列出所有记录
func (s *Storage) List() ([]*StoredRecord, error) {
	if err := s.load(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*StoredRecord
	for _, record := range s.records {
		result = append(result, record)
	}

	return result, nil
}

// Delete 删除记录
func (s *Storage) Delete(name string) error {
	if err := s.load(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.records[name]; !ok {
		return fmt.Errorf("record not found: %s", name)
	}

	delete(s.records, name)
	return s.save()
}
