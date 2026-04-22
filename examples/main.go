// examples/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/yourusername/smartdns/smartdns"
	"github.com/yourusername/smartdns/smartdns/provider/local"
)

func main() {
	// 创建临时目录用于演示
	tmpDir, err := os.MkdirTemp("", "smartdns_demo_*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	hostsPath := filepath.Join(tmpDir, "hosts")
	storagePath := filepath.Join(tmpDir, "records.json")

	fmt.Println("=== SmartDNS 组件库使用示例 ===\n")

	// 示例 1: 直接使用本地 Provider
	fmt.Println("1. 创建本地 Provider")
	provider := local.NewProvider(&local.Config{
		HostsPath:   hostsPath,
		StoragePath: storagePath,
		Domain:      "smartdns.local",
	})

	// 创建记录 - 支持多种 IP 格式
	fmt.Println("\n2. 创建 DNS 记录（支持多种IP格式）")
	record, err := provider.Create("app", []string{
		"1.1.1.1",         // 单IP
		"1.1.1.5-1.1.1.7", // IP范围
		"2.2.2.0/30",      // CIDR段 (2.2.2.0 - 2.2.2.3)
	})
	if err != nil {
		log.Fatalf("创建记录失败: %v", err)
	}
	fmt.Printf("   CNAME: %s\n", record.CNAME)
	fmt.Printf("   IPs: %v\n", record.IPs)

	// 查询记录
	fmt.Println("\n3. 查询 DNS 记录")
	record, err = provider.Get("app")
	if err != nil {
		log.Fatalf("查询记录失败: %v", err)
	}
	fmt.Printf("   Name: %s\n", record.Name)
	fmt.Printf("   CNAME: %s\n", record.CNAME)

	// 更新记录
	fmt.Println("\n4. 更新 DNS 记录")
	record, err = provider.Update("app", []string{"192.168.1.1", "192.168.1.2"})
	if err != nil {
		log.Fatalf("更新记录失败: %v", err)
	}
	fmt.Printf("   更新后 IPs: %v\n", record.IPs)

	// 列出所有记录
	fmt.Println("\n5. 列出所有记录")
	records, err := provider.List()
	if err != nil {
		log.Fatalf("列出记录失败: %v", err)
	}
	fmt.Printf("   总记录数: %d\n", len(records))
	for _, r := range records {
		fmt.Printf("   - %s -> %v\n", r.Name, r.CNAME)
	}

	// 示例 2: 创建另一个应用
	fmt.Println("\n6. 创建第二个应用")
	record2, err := provider.Create("api", []string{
		"10.0.0.1-10.0.0.5", // IP 范围
	})
	if err != nil {
		log.Fatalf("创建记录失败: %v", err)
	}
	fmt.Printf("   CNAME: %s\n", record2.CNAME)

	// 再次列出
	records, _ = provider.List()
	fmt.Printf("   总记录数: %d\n", len(records))

	// 示例 3: 使用工厂函数（需要配置文件）
	fmt.Println("\n7. 使用工厂函数创建 Provider")
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := fmt.Sprintf(`hosts_path: "%s"
storage_path: "%s"
domain: "demo.local"`, hostsPath, storagePath)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		log.Fatalf("写入配置失败: %v", err)
	}

	dns, err := smartdns.New("local", configPath)
	if err != nil {
		log.Fatalf("创建 Provider 失败: %v", err)
	}

	record3, err := dns.Create("web", []string{"172.16.0.1"})
	if err != nil {
		log.Fatalf("创建记录失败: %v", err)
	}
	fmt.Printf("   通过工厂创建 CNAME: %s\n", record3.CNAME)

	// 清理: 删除示例记录
	fmt.Println("\n8. 清理示例记录")
	provider.Delete("app")
	provider.Delete("api")
	dns.Delete("web")
	fmt.Println("   已删除所有测试记录")

	fmt.Println("\n=== 示例完成 ===")
	fmt.Printf("\n临时文件位置: %s\n", tmpDir)
	fmt.Println("注意: hosts 文件内容可在实际使用时查看")
}
