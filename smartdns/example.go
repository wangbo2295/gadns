// smartdns/example.go
package smartdns

import (
	"fmt"

	"github.com/yourusername/smartdns/smartdns/provider/local"
)

func ExampleLocalProvider() {
	// 创建本地Provider用于管理DNS记录
	provider := local.NewProvider(&local.Config{
		HostsPath:   "/tmp/hosts",
		StoragePath: "/tmp/records.json",
		Domain:      "smartdns.local",
	})

	// 创建记录 - 支持多种IP格式
	record, err := provider.Create("app", []string{
		"1.1.1.1",         // 单IP
		"1.1.1.5-1.1.1.7", // IP范围
		"2.2.2.0/30",      // CIDR段
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Created: %s\n", record.CNAME)

	// 查询记录
	record, _ = provider.Get("app")
	fmt.Printf("Record: %s -> %v\n", record.CNAME, record.IPs)

	// 列出所有记录
	records, _ := provider.List()
	fmt.Printf("Total: %d records\n", len(records))

	// Output:
	// Created: app.smartdns.local
	// Record: app.smartdns.local -> [1.1.1.1 1.1.1.5-1.1.1.7 2.2.2.0/30]
	// Total: 1 records
}

func ExampleLocalProvider_update() {
	provider := local.NewProvider(&local.Config{
		HostsPath:   "/tmp/hosts",
		StoragePath: "/tmp/records.json",
		Domain:      "smartdns.local",
	})

	// 创建记录
	provider.Create("web", []string{"192.168.1.1"})

	// 更新记录
	record, _ := provider.Update("web", []string{"192.168.1.10", "192.168.1.11"})
	fmt.Printf("Updated: %s -> %v\n", record.CNAME, record.IPs)

	// 删除记录
	provider.Delete("web")
	fmt.Println("Deleted")

	// Output:
	// Updated: web.smartdns.local -> [192.168.1.10 192.168.1.11]
	// Deleted
}

func ExampleSmartDNS_interface() {
	// SmartDNS 接口的完整用法示例

	// 创建本地provider实现SmartDNS接口
	var dns SmartDNS
	dns = local.NewProvider(&local.Config{
		HostsPath:   "/tmp/hosts",
		StoragePath: "/tmp/records.json",
		Domain:      "example.local",
	})

	// CRUD 操作
	record, _ := dns.Create("db", []string{"10.0.1.1", "10.0.1.2"})
	fmt.Printf("Create: %s\n", record.CNAME)

	record, _ = dns.Get("db")
	fmt.Printf("Get: %s\n", record.Name)

	record, _ = dns.Update("db", []string{"10.0.2.1"})
	fmt.Printf("Update: %v\n", record.IPs)

	dns.Delete("db")
	fmt.Println("Delete: success")

	// Output:
	// Create: db.example.local
	// Get: db
	// Update: [10.0.2.1]
	// Delete: success
}

func ExampleIPFormats() {
	// 展示支持的IP格式

	fmt.Println("支持的IP格式:")
	fmt.Println("单IP: 1.1.1.1")
	fmt.Println("IP范围: 1.1.1.1-1.1.1.10")
	fmt.Println("CIDR段: 1.1.1.0/24")
	fmt.Println("混合: [\"1.1.1.1\", \"1.1.1.5-1.1.1.10\", \"2.2.2.0/30\"]")

	// Output:
	// 支持的IP格式:
	// 单IP: 1.1.1.1
	// IP范围: 1.1.1.1-1.1.1.10
	// CIDR段: 1.1.1.0/24
	// 混合: ["1.1.1.1", "1.1.1.5-1.1.1.10", "2.2.2.0/30"]
}
