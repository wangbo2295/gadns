// examples/tencent/main.go - 通过 factory 创建腾讯云 Provider 的完整示例
package main

import (
	"fmt"
	"os"

	"github.com/wangbo2295/gadns/config"
	"github.com/wangbo2295/gadns/provider"
)

func main() {
	configPath := os.Getenv("HOME") + "/.gadns/tencent.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	// 读取域名，用于构造完整域名
	cfg, err := config.Load[config.TencentConfig](configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	cp, err := provider.New("tencent", configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建 provider 失败: %v\n", err)
		os.Exit(1)
	}

	// 使用配置中的域名构造完整域名
	fqdn := func(name string) string { return name + "." + cfg.Domain }

	// 1. 创建记录
	fmt.Println("=== 1. 创建记录 ===")

	record, err := cp.Create(fqdn("app"), []string{"1.1.1.1"})
	mustOK(err, "创建 app")
	fmt.Printf("  %s -> cname=%s, ips=%v\n", fqdn("app"), record.CNAME, record.IPs)

	record, err = cp.Create(fqdn("web"), []string{"10.0.0.1", "10.0.0.2"})
	mustOK(err, "创建 web")
	fmt.Printf("  %s -> cname=%s, ips=%v\n", fqdn("web"), record.CNAME, record.IPs)

	// 2. 查询
	fmt.Println("\n=== 2. 查询记录 ===")
	r, err := cp.Get(fqdn("app"))
	mustOK(err, "查询 app")
	fmt.Printf("  %s -> cname=%s, ips=%v\n", fqdn("app"), r.CNAME, r.IPs)

	// 3. 列出
	fmt.Println("\n=== 3. 列出所有记录 ===")
	records, err := cp.List()
	mustOK(err, "列出全部")
	for _, rec := range records {
		fmt.Printf("  %s -> cname=%s, ips=%v\n", rec.Name, rec.CNAME, rec.IPs)
	}

	// 4. 更新
	fmt.Println("\n=== 4. 更新记录 ===")
	updated, err := cp.Update(fqdn("app"), []string{"2.2.2.2"})
	mustOK(err, "更新 app")
	fmt.Printf("  %s -> cname=%s, ips=%v\n", fqdn("app"), updated.CNAME, updated.IPs)

	// 5. 清理
	fmt.Println("\n=== 5. 清理 ===")
	err = cp.Delete(fqdn("app"))
	mustOK(err, "删除 app")
	err = cp.Delete(fqdn("web"))
	mustOK(err, "删除 web")

	fmt.Println("\n=== 完成 ===")
}

func mustOK(err error, action string) {
	if err != nil {
		fmt.Printf("  [失败] %s: %v\n", action, err)
		os.Exit(1)
	}
}
