// SmartDNS CLI - 主程序
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/smartdns/core"
	"github.com/yourusername/smartdns/provider"
)

// Version 信息（由 ldflags 设置）
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add", "update":
		runAddUpdate(command)
	case "get":
		runGet()
	case "list":
		runList()
	case "delete":
		runDelete()
	case "help", "-h", "--help":
		showHelp()
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n\n", command)
		showHelp()
		os.Exit(1)
	}
}

func runAddUpdate(command string) {
	fs := flag.NewFlagSet(command, flag.ExitOnError)
	ipsFlag := fs.String("ips", "", "IP addresses (comma-separated)")
	configFlag := fs.String("c", "", "Config file path")

	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: missing record name\nUsage: smartdns %s -ips <ips> <name>\n", command)
		os.Exit(1)
	}

	name := fs.Arg(0)
	ips := *ipsFlag
	configPath := *configFlag

	if ips == "" {
		fmt.Fprintln(os.Stderr, "Error: -ips flag is required")
		os.Exit(1)
	}

	// 解析 IP 列表
	ipList := strings.Split(ips, ",")

	// 创建 provider
	dns, err := provider.New("local", getConfigPath(configPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create provider: %v\n", err)
		os.Exit(1)
	}

	var record *core.Record
	if command == "add" {
		record, err = dns.Create(name, ipList)
	} else {
		record, err = dns.Update(name, ipList)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ %s successful\n", command)
	fmt.Printf("  Name:  %s\n", record.Name)
	fmt.Printf("  CNAME: %s\n", record.CNAME)
	fmt.Printf("  IPs:   %s\n", strings.Join(record.IPs, ", "))
}

func runGet() {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	configFlag := fs.String("c", "", "Config file path")

	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: missing record name\nUsage: smartdns get <name>")
		os.Exit(1)
	}

	name := fs.Arg(0)
	configPath := *configFlag

	dns, err := provider.New("local", getConfigPath(configPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create provider: %v\n", err)
		os.Exit(1)
	}

	record, err := dns.Get(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Name:   %s\n", record.Name)
	fmt.Printf("CNAME:  %s\n", record.CNAME)
	fmt.Printf("IPs:    %s\n", strings.Join(record.IPs, ", "))
}

func runList() {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	configFlag := fs.String("c", "", "Config file path")

	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	configPath := *configFlag

	dns, err := provider.New("local", getConfigPath(configPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create provider: %v\n", err)
		os.Exit(1)
	}

	records, err := dns.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(records) == 0 {
		fmt.Println("No records found")
		return
	}

	fmt.Printf("Total Records: %d\n\n", len(records))
	fmt.Printf("%-15s | %-25s | %s\n", "Name", "CNAME", "IPs")
	fmt.Printf("%s-+-%s-+-%s\n", strings.Repeat("-", 15), strings.Repeat("-", 25), strings.Repeat("-", 40))

	for _, r := range records {
		ips := strings.Join(r.IPs, ", ")
		if len(ips) > 37 {
			ips = ips[:37] + "..."
		}
		fmt.Printf("%-15s | %-25s | %s\n", r.Name, r.CNAME, ips)
	}
}

func runDelete() {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	configFlag := fs.String("c", "", "Config file path")

	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: missing record name\nUsage: smartdns delete <name>")
		os.Exit(1)
	}

	name := fs.Arg(0)
	configPath := *configFlag

	dns, err := provider.New("local", getConfigPath(configPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create provider: %v\n", err)
		os.Exit(1)
	}

	if err := dns.Delete(name); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Record '%s' deleted\n", name)
}

func getConfigPath(configPath string) string {
	if configPath != "" {
		return configPath
	}
	return "~/.smartdns/local.yaml"
}

func showHelp() {
	help := `SmartDNS CLI - DNS 记录管理工具

Version: ` + Version + ` (commit: ` + GitCommit + `)

Usage:
  smartdns <command> [options] [arguments]

Commands:
  add       新增 DNS 记录
  update    更新 DNS 记录
  get       查询 DNS 记录
  list      列出所有记录
  delete    删除 DNS 记录
  help      显示此帮助信息

IP 格式支持:
  单 IP:       1.1.1.1
  IP 范围:     1.1.1.1-1.1.1.10
  CIDR 段:     192.168.1.0/24
  混合格式:    1.1.1.1,10.0.0.1-10.0.0.5,192.168.1.0/30

Options:
  -ips <list>    IP 地址列表（逗号分隔）[add/update 必需]
  -c <path>      配置文件路径 (默认: ~/.smartdns/local.yaml)
  -h, --help     显示帮助信息

Examples:
  # 新增记录（单 IP）
  smartdns add -ips 1.1.1.1 app

  # 多个 IP
  smartdns add -ips 1.1.1.1,2.2.2.2 web

  # IP 范围
  smartdns add -ips 10.0.0.1-10.0.0.10 api

  # CIDR 段
  smartdns add -ips 192.168.1.0/24 db

  # 混合格式
  smartdns add -ips 1.1.1.1,10.0.0.1-10.0.0.5,192.168.1.0/30 cache

  # 使用自定义配置
  smartdns add -ips 1.1.1.1 -c /path/to/config.yaml app

  # 查询记录
  smartdns get app

  # 列出所有记录
  smartdns list

  # 更新记录
  smartdns update -ips 5.5.5.5 app

  # 删除记录
  smartdns delete app
`
	fmt.Println(help)
}
