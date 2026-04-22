// provider/local/hosts.go
package local

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// HostsManager hosts文件管理器
type HostsManager struct {
	path string
	mu   sync.Mutex
}

// NewHostsManager 创建hosts文件管理器
func NewHostsManager(path string) *HostsManager {
	return &HostsManager{
		path: path,
	}
}

// AddEntry 添加hosts条目
func (hm *HostsManager) AddEntry(ip, hostname string) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	// 读取现有内容
	lines, err := hm.readLines()
	if err != nil {
		return err
	}

	// 检查是否已存在
	for _, line := range lines {
		if hm.isEntryForHost(line, hostname) {
			// 已存在，更新
			return hm.updateEntryInLines(lines, ip, hostname)
		}
	}

	// 添加新条目
	lines = append(lines, fmt.Sprintf("%s\t%s", ip, hostname))
	return hm.writeLines(lines)
}

// UpdateEntry 更新hosts条目
func (hm *HostsManager) UpdateEntry(ip, hostname string) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	lines, err := hm.readLines()
	if err != nil {
		return err
	}

	return hm.updateEntryInLines(lines, ip, hostname)
}

// RemoveEntry 删除hosts条目
func (hm *HostsManager) RemoveEntry(hostname string) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	lines, err := hm.readLines()
	if err != nil {
		return err
	}

	var newLines []string
	for _, line := range lines {
		if !hm.isEntryForHost(line, hostname) {
			newLines = append(newLines, line)
		}
	}

	return hm.writeLines(newLines)
}

// readLines 读取hosts文件行
func (hm *HostsManager) readLines() ([]string, error) {
	file, err := os.Open(hm.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// writeLines 写入hosts文件行
func (hm *HostsManager) writeLines(lines []string) error {
	// 确保目录存在
	dir := filepath.Dir(hm.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(hm.path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range lines {
		if _, err := fmt.Fprintln(file, line); err != nil {
			return err
		}
	}

	return nil
}

// isEntryForHost 检查行是否是特定主机的条目
func (hm *HostsManager) isEntryForHost(line, hostname string) bool {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return false
	}

	fields := strings.Fields(line)
	for i, field := range fields {
		if i > 0 && field == hostname {
			return true
		}
	}
	return false
}

// updateEntryInLines 在行列表中更新条目
func (hm *HostsManager) updateEntryInLines(lines []string, ip, hostname string) error {
	for i, line := range lines {
		if hm.isEntryForHost(line, hostname) {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				fields[0] = ip
				lines[i] = strings.Join(fields, "\t")
			}
			return hm.writeLines(lines)
		}
	}
	return fmt.Errorf("entry not found for hostname: %s", hostname)
}
